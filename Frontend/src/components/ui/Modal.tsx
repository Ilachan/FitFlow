import React, { useEffect } from "react";
import { Icons } from "../../lib/icons";
import Card from "./Card";

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

const Modal = ({ isOpen, onClose, children }: ModalProps) => {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };

    if (isOpen) {
      document.body.style.overflow = "hidden";
      window.addEventListener("keydown", handleEscape);
    }

    return () => {
      document.body.style.overflow = "unset";
      window.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-100 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="fixed inset-0 cursor-pointer" onClick={onClose} />
      <Card className="max-w-md w-full p-0 overflow-hidden relative z-10 animate-in zoom-in-95 duration-200 shadow-2xl ring-1 ring-slate-200/50">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 p-2 rounded-full hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors z-20 bg-white/80 backdrop-blur"
        >
          <Icons.X className="w-5 h-5" />
        </button>

        {children}
      </Card>
    </div>
  );
};

export default Modal;
