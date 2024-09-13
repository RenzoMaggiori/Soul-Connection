import { Skeleton } from "@/components/ui/skeleton";
import React from "react";

interface LoadingSkeletonProps {
  isLoading: boolean;
  data: boolean;
  skeletonClassName?: string;
  children: React.ReactNode;
}

export function LoadingSkeleton({
  isLoading,
  data,
  skeletonClassName,
  children,
}: LoadingSkeletonProps) {
  if (isLoading || !data) {
    return <Skeleton className={skeletonClassName} />;
  }
  return <>{children}</>;
}
