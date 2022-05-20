export const LoadingSkeleton = () => {
  return (
    <div className="mt-12 mx-4">
      <div className="animate-pulse flex space-x-4 border-l border-slate-300 pl-6 mt-4 p-4 ml-4">
        <div className="flex-1 space-y-6 py-1">
          <div className="h-2 bg-slate-200 rounded"></div>
          <div className="space-y-3">
            <div className="grid grid-cols-3 gap-4">
              <div className="h-2 bg-slate-200 rounded col-span-2"></div>
              <div className="h-2 bg-slate-200 rounded col-span-1"></div>
            </div>
            <div className="h-2 bg-slate-200 rounded"></div>
            <div className="grid grid-cols-4 gap-4">
              <div className="h-2 bg-slate-200 rounded col-span-2"></div>
              <div className="h-2 bg-slate-200 rounded col-span-2"></div>
            </div>
            <div className="grid grid-cols-4 gap-4">
              <div className="h-2 bg-slate-200 rounded col-span-1"></div>
              <div className="h-2 bg-slate-200 rounded col-span-2"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
