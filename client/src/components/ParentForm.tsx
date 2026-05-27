import { ReactNode } from "react"
import { ModalCloseButton } from "./ModalCloseButton"
import '../styles/ParentForm.css'

type ModalProps = {
  children: ReactNode;
  isDialog:boolean;
  isOpen: boolean; 
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  contentClassName?: string;
  backdropClassName?: string;
};

export const ParentForm: React.FC<ModalProps> = ({
  children,
  isDialog,
  isOpen,
  setIsOpen,
  contentClassName,
  backdropClassName,
}) => {
    const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
        if (e.target === e.currentTarget) {
            setIsOpen(false);
        }
    };

    return (
        <div className={isOpen ? "parent-form-host parent-form-host--open" : "parent-form-host"}>
            {isOpen &&
            <div
                id="form"
                className={["background", backdropClassName].filter(Boolean).join(" ")}
                onClick={(e) => handleBackdropClick(e)}
            >
                <div className={["content", contentClassName].filter(Boolean).join(" ")}>
                    {isDialog && (
                        <ModalCloseButton onClick={() => setIsOpen(false)} />
                    )}
                    {children}
                </div>
            </div>}
        </div>
    )
}