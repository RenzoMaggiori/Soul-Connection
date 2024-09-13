import { Skeleton } from "./skeleton";

export function SkeletonTable() {
    const numColumns = 4;
    const numRows = 10;

    return (
      <div className="w-full my-12 border rounded-lg overflow-hidden">
        <div className="bg-gray-200 p-4">
          <Skeleton className="h-6 w-1/4" /> {/* Skeleton for table header */}
        </div>
        <div className="p-4">
          <table className="w-full">
            <thead>
              <tr>
                {Array.from({ length: numColumns }).map((_, index) => (
                  <th key={index} className="px-4 py-2">
                    <Skeleton className="h-4 w-full bg-slate-200" /> {/* Skeleton for column headers */}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {Array.from({ length: numRows }).map((_, rowIndex) => (
                <tr key={rowIndex}>
                  {Array.from({ length: numColumns }).map((_, colIndex) => (
                    <td key={colIndex} className="px-4 py-4">
                      <Skeleton className="h-6 w-full bg-slate-200" /> {/* Skeleton for table cells */}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    );
  }