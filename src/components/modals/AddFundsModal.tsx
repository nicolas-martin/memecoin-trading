import React from 'react';
import { Dialog } from '@headlessui/react';
import AddFundsForm from '../payment/AddFundsForm';

interface AddFundsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AddFundsModal: React.FC<AddFundsModalProps> = ({ isOpen, onClose }) => {
  const handleSuccess = () => {
    // Show success message
    onClose();
  };

  const handleError = (error: Error) => {
    // Show error message
    console.error('Payment failed:', error);
  };

  return (
    <Dialog open={isOpen} onClose={onClose} className="relative z-50">
      <div className="fixed inset-0 bg-black/30" aria-hidden="true" />
      <div className="fixed inset-0 flex items-center justify-center p-4">
        <Dialog.Panel className="mx-auto max-w-sm rounded-lg bg-white">
          <div className="p-6">
            <Dialog.Title as="h3" className="text-lg font-medium leading-6 text-gray-900 mb-4">
              Add Funds
            </Dialog.Title>
            <AddFundsForm
              onSuccess={handleSuccess}
              onError={handleError}
            />
          </div>
        </Dialog.Panel>
      </div>
    </Dialog>
  );
};

export default AddFundsModal; 