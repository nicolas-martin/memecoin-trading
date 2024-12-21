const LoadingSpinner = () => {
  return (
    <div className="min-h-screen bg-black flex items-center justify-center">
      <div className="flex flex-col items-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
        <span className="text-white mt-4">Loading...</span>
      </div>
    </div>
  );
};

export default LoadingSpinner; 