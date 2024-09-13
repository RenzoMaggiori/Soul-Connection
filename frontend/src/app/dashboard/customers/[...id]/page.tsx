"use client";

import React from "react";
import Image from "next/image";
import { useQuery } from "@tanstack/react-query";
import { ScrollArea } from "@radix-ui/react-scroll-area";
import { Mail, Bookmark } from "lucide-react";
import { useRouter } from "next/navigation";
import { Icons } from "@/components/icons";
import { DataTable } from "@/components/ui/data-table";
import { meetingColumns } from "./columnsMeeting";
import { paymentColumns } from "./columnsPayment"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { getCustomerById, getCustomerEmployeeById } from "@/db/customer";
import { getEncountersByCustomerId } from "@/db/encounter";
import { getPaymentsByCustomerId } from "@/db/payment";
import { getSession } from "@/db/db";


const isNumeric = (val: string): boolean => !isNaN(Number(val));

const getImage = async (customerId: number): Promise<Blob | null> => {
  const session = await getSession();
  if (!session) return null;

  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/customers/${customerId}/image`,
      {
        headers: {
          Authorization: `${session.token}`,
        },
      }
    );

    if (!response.ok) {
      console.error('Failed to fetch image:', response.statusText);
      return null;
    }

    return await response.blob();
  } catch (error) {
    console.error('Error fetching image:', error);
    return null;
  }
};


export default function ProfilePage({ params }: { params: { id: string } }) {
  const router = useRouter();
  const customerId = parseInt(params.id[0]);

  // TO DO: employee

  const { data, isLoading, isError } = useQuery({
    queryFn: async () => {
      const [customer, encounters, employee, payments, imageData] = await Promise.all([
        getCustomerById(customerId),
        getEncountersByCustomerId(customerId),
        getCustomerEmployeeById(customerId),
        getPaymentsByCustomerId(customerId),
        // getImagesByCustomerId(customerId),
        getImage(customerId),
      ]);
      return { customer, encounters, employee, payments, imageData };
    },
    queryKey: ["CustomerData", params.id[0]],
    gcTime: 1000 * 60,
  });

  if (!isNumeric(params.id[0])) {
    return <div>Invalid id</div>;
  }

  if (isError) return <div>Error...</div>;

  const imageUrl = window.URL.createObjectURL(data?.imageData ?? new Blob());

  const renderProfileHeader = () => (
    <div className="flex flex-col items-center justify-center gap-4 p-4 pb-0">
      <Avatar className="h-20 w-20">
        <AvatarImage src={imageUrl} alt="Profile Picture" />
        <AvatarFallback>
          {!isLoading && data?.customer ? (
            data.customer.Name[0] + data.customer.Surname[0]
          ) : (
            <Skeleton className="h-20 w-20 rounded-full" />
          )}
        </AvatarFallback>
      </Avatar>
      <h2 className="text-xl font-bold">
        {!isLoading && data?.customer ? (
          `${data.customer.Name} ${data.customer.Surname}`
        ) : (
          <Skeleton className="h-8 w-32" />
        )}
      </h2>
    </div>
  );

  const renderProfileStats = () => (
    <div className="flex justify-between p-4">
      <div className="flex flex-col items-center gap-2">
        <h2 className="text-xl font-bold">
          {!isLoading && data?.encounters ? (
            `${data?.encounters?.length}`
          ) : (
            <Skeleton className="h-5 w-10" />)}
        </h2>
        <div className="flex flex-col items-center">
          <p>Total</p>
          <p>Encounters</p>
        </div>
      </div>
      <div className="flex flex-col items-center gap-2">
        <h2 className="text-xl font-bold">
          {!isLoading && data?.encounters ? (
            `${data?.encounters?.filter((encounter) => encounter.Rating > 3).length}`
          ) : (
            <Skeleton className="h-5 w-10" />)}
        </h2>
        <p>Positives</p>
      </div>
      <div className="flex flex-col items-center gap-2">
        <h2 className="text-xl font-bold">
          {!isLoading && data?.encounters ? (
            `${data?.encounters?.filter((encounter) => encounter.Date < new Date().toISOString()).length}`
          ) : (
            <Skeleton className="h-5 w-10" />
          )}
        </h2>
        <p>In Progress</p>
      </div>
    </div>
  );

  const renderProfileDetails = () => (
    <div className="p-4 pb-0">
      <h3 className="text-lg mb-4">SHORT DETAILS</h3>
      <div className="space-y-4 text-sm">
        <div className="space-y-1">
          <p>User ID:</p>
          <span className="font-semibold">
            {!isLoading && data?.customer ? (
              `${data?.customer?.Soul_Connection_Id}`
            ) : (
              <Skeleton className="h-5 w-full" />
            )}
          </span>
        </div>
        <div className="space-y-1">
          <p>Email:</p>
          <span className="font-semibold">
            {!isLoading && data?.customer ? (
              `${data?.customer?.Email}`
            ) : (
              <Skeleton className="h-5 w-full" />
            )}
          </span>
        </div>
        <div className="space-y-1">
          <p>Address:</p>
          <span className="font-semibold">
            {!isLoading && data?.customer ? (
              `${data?.customer?.Address}`
            ) : (
              <Skeleton className="h-5 w-full" />
            )}
          </span>
        </div>
        <div className="space-y-1">
          <p>Astrological Sign:</p>
          <span className="font-semibold">
            {!isLoading && data?.customer ? (
              `${data?.customer?.Astrological_Sign}`
            ) : (
              <Skeleton className="h-5 w-full" />
            )}
          </span>
        </div>
      </div>
    </div>
  );

  const renderProfile = () => (
    <div className="lg:w-1/3 space-y-6 bg-white rounded-md shadow">
      {renderProfileHeader()}
      <Separator />
      <div className="flex justify-center space-x-6">
        <Mail className="h-5 w-5" />
        <Bookmark className="h-5 w-5" />
      </div>
      <Separator />
      {renderProfileStats()}
      <Separator />
      {renderProfileDetails()}
    </div>
  );

  const renderMeetings = () => (
    <div className="space-y-4">
      <h3 className="mb-2 text-lg font-semibold">Recent Meetings</h3>
      {!isLoading && data?.encounters ? (
        <DataTable
          data={data?.encounters}
          columns={meetingColumns}
          pageSize={5}
        />
      ) : (
        <Skeleton className="w-full h-25" />
      )}
    </div>
  );

  const renderPaymentHistory = () => (
    <div className="space-y-4">
      <h3 className="mb-2 text-lg font-semibold">Payments History</h3>
      {!isLoading && data?.payments ? (
        <DataTable
          data={data?.payments}
          columns={paymentColumns}
          pageSize={5}
        />
      ) : (
        <Skeleton className="w-full h-25" />
      )}
    </div>
  );

  return (
    <div className="mx-auto px-4 py-4 sm:px-6 lg:px-8 text-generic">
      <div className="flex justify-between">
        <h1 className="text-xl mb-6 font-bold md:text-2xl">Customer Details</h1>
        <Button variant={"generic"} onClick={() => router.push("/dashboard/customers")}>
          Back
        </Button>
      </div>
      <div className="flex flex-col lg:flex-row gap-6 bg-gray-100 min-h-screen">
        {renderProfile()}
        <div className="lg:w-3/4 space-y-6 bg-white rounded-md shadow p-4">
          {renderMeetings()}
          {renderPaymentHistory()}
        </div>
      </div>
    </div>
  );
}
