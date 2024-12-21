const Leaderboard = () => {
  return (
    <div className="min-h-screen bg-gray-100 pb-20">
      <div className="bg-white">
        <div className="max-w-md mx-auto px-4 py-4">
          <h1 className="text-2xl font-semibold">Leaderboard</h1>
        </div>
      </div>

      <div className="max-w-md mx-auto px-4 py-4">
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
          <div className="px-4 py-3 border-b border-gray-100">
            <div className="flex justify-between items-center">
              <span className="text-sm font-medium text-gray-500">Rank</span>
              <span className="text-sm font-medium text-gray-500">Profit</span>
            </div>
          </div>
          
          <div className="flex justify-center items-center h-32 text-gray-400">
            No traders yet
          </div>
        </div>
      </div>
    </div>
  );
};

export default Leaderboard; 