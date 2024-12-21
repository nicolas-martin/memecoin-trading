import React, { useState } from 'react';
import AddFundsModal from '../components/modals/AddFundsModal';

const Profile: React.FC = () => {
  const [showAddFunds, setShowAddFunds] = useState(false);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="bg-white rounded-lg shadow-lg p-6">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Profile</h1>
          <button
            onClick={() => setShowAddFunds(true)}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
          >
            Add Funds
          </button>
        </div>

        {/* Other profile content */}
      </div>

      <AddFundsModal
        isOpen={showAddFunds}
        onClose={() => setShowAddFunds(false)}
      />
    </div>
  );
};

export default Profile; 