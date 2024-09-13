"use client";

import React, { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import UserAvatar from "./UserAvatar";
import { LinearProgressBar } from "react-percentage-bar";
import { Skeleton } from "@/components/ui/skeleton";
import { getCustomers } from "@/db/customer";

export default function AstrologicalCompatibility() {
  const [user1, setUser1] = useState<string | undefined>(undefined);
  const [user2, setUser2] = useState<string | undefined>(undefined);
  const [percentage, setPercentage] = useState<number>(0);
  const { isLoading, isError, data } = useQuery({
    queryFn: async () => {
      const customers = await getCustomers();
      return { customers };
    },
    queryKey: ["CustomerData"],
    gcTime: 1000 * 60,
  });

  useEffect(() => {
    if (user1 && user2 && user1 != user2) {
      const randomPercentage = Math.floor(Math.random() * 100);
      setPercentage(randomPercentage + 1);
    } else {
      setPercentage(0);
    }
  }, [user1, user2]);

  if (isError) return <div>Error...</div>;
  if (isLoading || !data) return <div>Loading...</div>;
  if (!data.customers) return <div>No data...</div>;
  return (
    <div className="h-full flex justify-center align-middle items-center">
    <div className="flex h-fit flex-col align-middle items-center justify-center bg-background p-4 md:p-6">
      <h1 className="mb-12 text-center text-generic text-2xl font-bold max-w-xl">
        Select two users to compare their astrological compatibility.
      </h1>
      <div className="grid w-full max-w-md grid-cols-2 gap-8">
        <UserAvatar
          user={user1}
          onSelectUser={setUser1}
          data={data.customers}
          isLoading={isLoading}
          selectedUser={user1 || ""}
          title="Select User 1"
        />
        <UserAvatar
          user={user2}
          onSelectUser={setUser2}
          data={data.customers}
          isLoading={isLoading}
          selectedUser={user2 || ""}
          title="Select User 2"
        />
      </div>
      <div key={percentage} className="pt-7">
        <LinearProgressBar percentage={percentage} />
      </div>
      <div className="mt-12 w-full max-w-md rounded-lg bg-muted p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-generic">Compatibility</h2>
        </div>
        <div className="mt-2 text-muted-foreground">
          {user1 && user2 ? (
            `Based on their profiles, ${user1} and ${user2} have a ${percentage}% compatibility match.`
          ) : (
            <Skeleton className="h-4 w-64 bg-zinc-300" />
          )}
        </div>
      </div>
    </div>
    </div>
  );
}
