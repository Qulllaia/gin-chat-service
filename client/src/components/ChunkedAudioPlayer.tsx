import React, { useState, useRef, useEffect } from 'react';

interface ChunkedAudioPlayerProps {
  ws: React.RefObject<WebSocket | null>
}

interface BufferLowEventDetail {
  loadedChunks: number;
  totalChunks: number;
  remaining: number;
}

type BufferStatus = 'stable' | 'low';

const ChunkedAudioPlayer = ({ ws }: ChunkedAudioPlayerProps) => {
  const [bufferStatus, setBufferStatus] = useState<BufferStatus>('stable');
  const [loadedPercentage, setLoadedPercentage] = useState<number>(0);
  const audioRef = useRef<HTMLAudioElement>(null);
  const mediaSourceRef = useRef<MediaSource | null>(null);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const mediaSource = new MediaSource();
    mediaSourceRef.current = mediaSource;

    audio.src = URL.createObjectURL(mediaSource);

    const handleSourceOpen = async () => {
      const sourceBuffer = mediaSource.addSourceBuffer('audio/mpeg');

      mediaSource.addEventListener('updateend', function() {
        if (!sourceBuffer.updating && mediaSource.readyState === 'open') {

          let buffered = audio.buffered;
          let currentTime = audio.currentTime
          if (buffered.length > 0) {
            let bufferedEnd = buffered.end(buffered.length - 1);
            let timeToEnd = bufferedEnd - currentTime;

            if (timeToEnd < 2) {
              loadNextChunk()
            }
          }
        }
      })

      if (ws.current) {
        ws.current.send(JSON.stringify({
          type: 'MEDIA',
          index: 1,
        }));

        ws.current.onmessage = (event) => {
          const message = JSON.parse(event.data);

          if (message.type === 'audio_chunk') {
            const binaryString = atob(message.data);
            const bytes = new Uint8Array(binaryString.length);

            for (let i = 0; i < binaryString.length; i++) {
              bytes[i] = binaryString.charCodeAt(i);
            }
            sourceBuffer.appendBuffer(bytes);
          }
        }

      }
    };

    mediaSource.addEventListener('sourceopen', handleSourceOpen);

    const handleBufferLow = (e: Event) => {
      const customEvent = e as CustomEvent<BufferLowEventDetail>;
      setBufferStatus('low');
      console.log('Буфер заканчивается!', customEvent.detail);
    };

    const witingHanlder = (e: Event) => {
      let buffered = audio.buffered;
      let currentTime = audio.currentTime
      if (buffered.length > 0) {
        let bufferedEnd = buffered.end(buffered.length - 1);
        let timeToEnd = bufferedEnd - currentTime;

        if (timeToEnd < 2) {
          loadNextChunk()
        }
      }
    }

    const loadNextChunk = () => {

      if (ws.current) {
        ws.current.send(JSON.stringify({
          type: 'MEDIA',
          index: 2,
        }));
      }
    }


    audio.addEventListener('bufferlow', handleBufferLow as EventListener);


    audio.addEventListener('timeupdate', witingHanlder as EventListener);

    return () => {
      mediaSource.removeEventListener('sourceopen', handleSourceOpen);
      audio.removeEventListener('bufferlow', handleBufferLow as EventListener);

      if (mediaSource.readyState === 'open') {
        mediaSource.endOfStream();
      }

      if (audio.src) {
        URL.revokeObjectURL(audio.src);
      }
    };
  }, []);


  return (
    <div>
      <audio ref={audioRef} controls />
      <div className={`buffer-status ${bufferStatus}`}>
        {bufferStatus === 'low' && (
          <div>Загружаем следующие чанки...</div>
        )}
        <progress value={loadedPercentage} max="100" />
      </div>
    </div>
  );
};

export default ChunkedAudioPlayer;
