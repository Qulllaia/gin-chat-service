import React, { useState, useRef, useEffect } from 'react';

interface ChunkedAudioPlayerProps {
  ws: React.RefObject<WebSocket | null>
}

const STATIC_BUFFER_READ_LIMIT = 1024 * 100
const ChunkedAudioPlayer = ({ ws }: ChunkedAudioPlayerProps) => {
  console.log('ChunkedAudioPlayer render, ws.current:', ws.current);

  const audioRef = useRef<HTMLAudioElement>(null);
  const mediaSourceRef = useRef<MediaSource | null>(null);
  const currentIndex = useRef<number>(0);

  const [currentIndexState, setCurrentIndexState] = useState<number>(0);
  const sourceBufferRef = useRef<SourceBuffer | null>(null);
  const [maxValue, setMaxValue] = useState<number>(0);
  const [maxBytes, setMaxBytes] = useState<number>(0);

  const [audioProgress, setAudioProgress] = useState<number>(0);

  useEffect(() => {
    console.log('useEffect запущен');

    const audio = audioRef.current;
    if (!audio) {
      console.log('audio элемент не найден');
      return;
    }

    const mediaSource = new MediaSource();
    mediaSourceRef.current = mediaSource;
    audio.src = URL.createObjectURL(mediaSource);

    console.log('MediaSource создан, URL:', audio.src);

    const handleMessage = (event: MessageEvent) => {
      const message = JSON.parse(event.data);
      console.log('Получено сообщение:', message.type);

      if (message.type === 'audio_chunk' && sourceBufferRef.current) {
        console.log('Добавляем аудио-чанк, индекс:', currentIndexState);

        if (currentIndexState === 0) {
          setMaxValue(message.size);
          setMaxBytes(message.size_bytes)
        }

        console.log('Максимальное значение:', maxValue);

        const binaryString = atob(message.data);

        const bytes = new Uint8Array(binaryString.length);

        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }

        if (!sourceBufferRef.current.updating) {
          sourceBufferRef.current.appendBuffer(bytes);
        } else {
          console.warn('SourceBuffer занят, пропускаем чанк');
        }
      }
    };

    const handleSourceOpen = async () => {
      console.log('sourceopen событие');

      if (!mediaSource.sourceBuffers.length) {
        const sourceBuffer = mediaSource.addSourceBuffer('audio/mpeg');
        sourceBufferRef.current = sourceBuffer;

        console.log('SourceBuffer создан');
      }

      if (ws.current) {
        ws.current.addEventListener('message', handleMessage);

        ws.current.send(JSON.stringify({
          type: 'MEDIA',
          index: currentIndexState,
        }));

        console.log('Отправлен запрос на чанк:', currentIndexState);
      }
    };

    const checkBuffer = () => {
      if (!audio) return;

      const buffered = audio.buffered;
      if (buffered.length === 0) return;

      const bufferedEnd = buffered.end(buffered.length - 1);


      const timeToEnd = bufferedEnd - audio.currentTime;
      setAudioProgress(audio.currentTime)
      console.log('Проверка буфера:', {
        bufferedEnd,
        currentTime: audio.currentTime,
        timeToEnd
      });

      if (timeToEnd < 2) {
        setCurrentIndexState(currentIndexState + 1);

        if (ws.current) {
          ws.current.send(JSON.stringify({
            type: 'MEDIA',
            index: currentIndexState,
          }));

          console.log('Запрошен следующий чанк:', currentIndexState);
        }
      }
    };

    mediaSource.addEventListener('sourceopen', handleSourceOpen);
    audio.addEventListener('timeupdate', checkBuffer);

    return () => {
      console.log('Очистка эффекта');

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
  }, [ws, maxValue, currentIndexState]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
    let value = Number(event.target.value)
    setAudioProgress(value)
    let bytesFromStart = value / maxValue * maxBytes;
    let newIndex = bytesFromStart / STATIC_BUFFER_READ_LIMIT;
    let rsultValue = parseInt(newIndex.toString())
    if (currentIndexState !== rsultValue) {
      setCurrentIndexState(rsultValue);
    }
  };
  return (
    <div>
      <audio ref={audioRef} controls />
      <label htmlFor="customRange1" className="form-label">Example range</label>
      <input onChange={handleChange} max={maxValue} value={audioProgress} type="range" className="form-range" id="customRange1" />

    </div >
  );
};

export default ChunkedAudioPlayer;
