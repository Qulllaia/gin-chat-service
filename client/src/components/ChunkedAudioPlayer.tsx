import React, { useState, useRef, useEffect, useCallback } from 'react';

interface ChunkedAudioPlayerProps {
  ws: React.RefObject<WebSocket | null>
}

const STATIC_BUFFER_READ_LIMIT = 1024 * 100
const ChunkedAudioPlayer = ({ ws }: ChunkedAudioPlayerProps) => {

  const audioRef = useRef<HTMLAudioElement>(null);
  const mediaSourceRef = useRef<MediaSource | null>(null);
  const currentIndex = useRef<number>(0);

  const sourceBufferRef = useRef<SourceBuffer | null>(null);
  const [durationValue, setDurationValue] = useState<number>(0);
  const [maxBytes, setMaxBytes] = useState<number>(0);

  const [audioProgress, setAudioProgress] = useState<number>(0);

  const isSelectedTimecode = useRef<boolean>(false);

  useEffect(() => {

    const audio = audioRef.current;
    if (!audio) {
      return;
    }

    const mediaSource = new MediaSource();
    mediaSourceRef.current = mediaSource;
    audio.src = URL.createObjectURL(mediaSource);


    mediaSource.addEventListener('sourceopen', handleSourceOpen);
    audio.addEventListener('timeupdate', checkBuffer);

    return () => {

      mediaSource.removeEventListener('sourceopen', handleSourceOpen);
      audio.removeEventListener('timeupdate', checkBuffer);

      if (ws.current) {
        ws.current.removeEventListener('message', handleMessage);
      }

      if (mediaSource.readyState === 'open') {
        mediaSource.endOfStream();
      }

      if (audio.src) {
        URL.revokeObjectURL(audio.src);
      }
    };
  }, [ws]);


  const handleMessage = async (event: MessageEvent) => {
    const message = JSON.parse(event.data);

    if (message.type === 'audio_end') {
      audioRef.current?.pause();
      setAudioProgress(0);
      currentIndex.current = 0;
    }

    if (message.type === 'audio_chunk' && sourceBufferRef.current) {
      console.warn("in")

      if (currentIndex.current === 0) {
        setDurationValue(message.audio_duration);
        setMaxBytes(message.size_bytes)
      }
      const binaryString = atob(message.data);

      const bytes = new Uint8Array(binaryString.length);

      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }

      if (isSelectedTimecode.current) {
        sourceBufferRef.current.abort();

        try {

          console.warn("sourceBufferRef.current.updating", sourceBufferRef.current.updating)
          if (sourceBufferRef.current.updating) {
            await waitForUpdateEnd(sourceBufferRef.current);
          }

          if (sourceBufferRef.current.buffered.length > 0) {
            const end = sourceBufferRef.current.buffered.end(sourceBufferRef.current.buffered.length - 1);
            sourceBufferRef.current.remove(0, end);
            await waitForUpdateEnd(sourceBufferRef.current);
          }

          console.warn("sourceBufferRef.current.timestampOffset ", sourceBufferRef.current.timestampOffset)
          sourceBufferRef.current.timestampOffset = 0;
          audioRef.current!.currentTime = 0;

          console.warn("sourceBufferRef.current.timestampOffset ", sourceBufferRef.current.timestampOffset)
          isSelectedTimecode.current = false;

          if (!sourceBufferRef.current.updating) {
            sourceBufferRef.current.appendBuffer(bytes);
          } else {
            setTimeout(() => {
              if (sourceBufferRef.current && !sourceBufferRef.current.updating) {
                sourceBufferRef.current.appendBuffer(bytes);
              }
            }, 100);
          }
        } catch (error) {
          console.error('Ошибка при очистке/добавлении буфера:', error);
        }


      } else {

        if (!sourceBufferRef.current.updating) {
          sourceBufferRef.current.appendBuffer(bytes);
        } else {
          setTimeout(() => {
            if (sourceBufferRef.current && !sourceBufferRef.current.updating) {
              sourceBufferRef.current.appendBuffer(bytes);
            }
          }, 100);
        }
      }

    }
  };

  const waitForUpdateEnd = (sourceBuffer: SourceBuffer): Promise<void> => {
    return new Promise((resolve, reject) => {
      const onUpdateEnd = () => {
        sourceBuffer.removeEventListener('updateend', onUpdateEnd);
        sourceBuffer.removeEventListener('error', onError);
        resolve();
      };

      const onError = (e: Event) => {
        sourceBuffer.removeEventListener('updateend', onUpdateEnd);
        sourceBuffer.removeEventListener('error', onError);
        reject(new Error('SourceBuffer update failed'));
      };

      sourceBuffer.addEventListener('updateend', onUpdateEnd, { once: true });
      sourceBuffer.addEventListener('error', onError, { once: true });
    });
  };

  const handleSourceOpen = async () => {
    if (mediaSourceRef.current) {

      if (!mediaSourceRef.current.sourceBuffers.length) {
        const sourceBuffer = mediaSourceRef.current.addSourceBuffer('audio/mpeg');
        sourceBufferRef.current = sourceBuffer;
      }

      if (ws.current) {
        if (ws.current.readyState !== WebSocket.CONNECTING) {

          ws.current.addEventListener('message', handleMessage);

          ws.current.send(JSON.stringify({
            type: 'MEDIA',
            index: currentIndex.current,
          }));
        }
      }
    }
  };

  const checkBuffer = () => {

    const audio = audioRef.current;
    if (!audio) return;

    const buffered = audio.buffered;
    if (buffered.length === 0) return;

    const bufferedEnd = buffered.end(buffered.length - 1);
    console.log(buffered.start(buffered.length - 1))


    // TODO: проблема перемотки в этой строчке. необходимо устранить проблему с тем, что
    // аудио не соответствует текущему буферу или сделать так, чтобы буфер считал конец не
    // завися от аудио
    const timeToEnd = bufferedEnd - audio.currentTime;
    setAudioProgress(STATIC_BUFFER_READ_LIMIT * currentIndex.current)

    if (timeToEnd < 2 && !sourceBufferRef.current?.updating && !isSelectedTimecode.current) {
      console.warn('pidars kotoriy vinoven v tom, chto u menya nichego ne rabotaet', currentIndex.current,
        bufferedEnd,
        buffered.start(buffered.length - 1),
        audio.currentTime)

      currentIndex.current = currentIndex.current + 1;

      if (ws.current) {
        ws.current.send(JSON.stringify({
          type: 'MEDIA',
          index: currentIndex.current,
        }));

      }
    }
  };

  const handleChangeSlider = (event: React.ChangeEvent<HTMLInputElement>, ws: React.RefObject<WebSocket | null>): void => {
    let value = Number(event.target.value)
    setAudioProgress(value)
    let index = value / STATIC_BUFFER_READ_LIMIT;
    let rsultValue = parseInt(index.toString())
    console.log(value, index, rsultValue)

    if (currentIndex.current !== rsultValue) {
      currentIndex.current = rsultValue;
      if (ws.current) {
        ws.current.send(JSON.stringify({
          type: 'MEDIA',
          index: currentIndex.current,
        }));
      }
    }

    isSelectedTimecode.current = true;
    console.warn('handleChangeSlider', isSelectedTimecode.current);
  };

  return (
    <div>
      <audio ref={audioRef} controls />
      <button onClick={() => {
        if (audioRef.current?.paused) {
          audioRef.current?.play();
        } else {
          audioRef.current?.pause();
        }
      }}>play</button>
      <input onChange={(event) => handleChangeSlider(event, ws)} max={maxBytes} value={audioProgress} type="range" className="form-range" id="customRange1" />

    </div >
  );
};

export default ChunkedAudioPlayer;
