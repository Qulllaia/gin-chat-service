import '../styles/ModalCloseButton.css';

type ModalCloseButtonProps = {
  onClick: () => void;
  className?: string;
};

export function ModalCloseButton({ onClick, className }: ModalCloseButtonProps) {
  return (
    <button
      type="button"
      className={['modal-close-btn', className].filter(Boolean).join(' ')}
      onClick={onClick}
      aria-label="Закрыть"
    >
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  );
}

export default ModalCloseButton;
