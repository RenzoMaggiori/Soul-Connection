"use client";

import React, { useState, useEffect } from "react";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import WardrobeItem from "@/components/ui/wardrobeItem";
import WardrobeOutfit from "@/components/ui/wardrobeOutfit";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { useQuery } from "@tanstack/react-query";
import { getCustomers } from "@/db/customer";
import { Customer } from "@/db/schemas";
import { Info, User } from "lucide-react";
import { getClothesByCustomerId } from "@/db/clothes";
import { getSession } from "@/db/db";

interface WardrobeItemType {
    id: number;
    image: string;
    type: string;
}

const WardrobePage = () => {
    const [items, setItems] = useState<WardrobeItemType[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedUser, setSelectedUser] = useState<string | undefined>(undefined);
    const [clothes, setClothes] = useState<WardrobeItemType[]>([]);

    const { isLoading, isError, data } = useQuery({
        queryFn: async () => {
            const customers = await getCustomers();
            return { customers };
        },
        queryKey: ["CustomerData"],
        gcTime: 1000 * 60,
    });

    const getImage = async (imageId: number) => {
        const session = await getSession();
        if (!session) return null;
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/clothes/${imageId}/image`, {
            headers: {
                Authorization: `${session.token}`,
            },
        });
        if (!response.ok) return null;
        return await response.blob();
    };

    useEffect(() => {
        const fetchClothesImages = async () => {
            const clothesWithImages = await Promise.all(clothes.map(async (clothe) => {
                const imageBlob = await getImage(clothe.id);
                const imageUrl = imageBlob ? URL.createObjectURL(imageBlob) : "/wardrobe/top1.jpeg";
                return {
                    id: clothe.id,
                    image: imageUrl,
                    type: clothe.type,
                };
            }));
            setItems(clothesWithImages);
            setLoading(false);
        };

        if (clothes.length > 0) {
            fetchClothesImages();
        }
    }, [clothes]);

    const handleSelectUser = async (customerName: string, customerId: number) => {
        setSelectedUser(customerName);
        setLoading(true);
        const clothesData = await getClothesByCustomerId(customerId);
        console.log("LOG clothesData", clothesData);

        const mappedClothes = clothesData?.map((clothe) => ({
            id: clothe.Id,
            image: clothe.Image_Id,
            type: clothe.Type,
        }));

        if (!mappedClothes) {
            setLoading(false);
            return;
        }
        setClothes(mappedClothes);
        setLoading(false);
    };

    const types = ["hat/cap", "top", "bottom", "shoes"];
    const groupedItems = types.map((type) => ({
        type,
        items: items.filter((item) => item.type === type),
    }));

    if (data && !data.customers) {
        return (
            <div className="flex h-screen">
                <div className="w-1/2 p-4">
                    <Skeleton />
                </div>
                <div className="w-1/2 p-4">
                    <Skeleton />
                </div>
            </div>
        );
    }

    return (
        <DndProvider backend={HTML5Backend}>
            <div className="flex h-full">
                <div className="w-1/2 space-y-4 p-4 overflow-y-auto">
                    <Dialog>
                        <DialogTrigger asChild>
                            <Button variant="outline">Choose a customer</Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-[425px]">
                            <DialogHeader>
                                <DialogTitle>Choose a customer</DialogTitle>
                            </DialogHeader>
                            <div className="flex items-center space-x-2">
                                <div className="grid flex-1 gap-2 text-generic">
                                    <Select
                                        onValueChange={(value) => {
                                            const customer = data?.customers?.find((cust: Customer) =>
                                                `${cust.Name} ${cust.Surname}` === value);
                                            if (customer) {
                                                handleSelectUser(value, customer.Id);
                                            }
                                        }}
                                        defaultValue={selectedUser}
                                    >
                                        <SelectTrigger className="w-[180px]">
                                            <SelectValue placeholder="Select a customer" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>Customer</SelectLabel>
                                                {data && data.customers?.map((customer: Customer) => (
                                                    <SelectItem
                                                        value={`${customer.Name} ${customer.Surname}`}
                                                        key={customer.Id}
                                                    >
                                                        {`${customer.Name} ${customer.Surname}`}
                                                    </SelectItem>
                                                ))}
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <DialogFooter className="sm:justify-start">
                                <DialogClose asChild>
                                    <Button type="button" variant="secondary">
                                        Close
                                    </Button>
                                </DialogClose>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>
                    {selectedUser ? (
                        <div className="flex items-center text-generic mb-2">
                            <User className="mr-2" size={22} />
                            <h2 className="text-2xl font-semibold text-generic">{selectedUser}</h2>
                        </div>
                    ) : (
                        <div className="flex items-center text-generic mb-2">
                            <Info className="mr-2" size={18} />
                            <p>There&rsquo;s no customer selected</p>
                        </div>
                    )}
                    {groupedItems.map((group) => (
                        <div key={group.type} className="mb-4">
                            <h2 className="mb-2 text-xl capitalize text-generic">
                                {group.type}
                            </h2>
                            <ScrollArea className="w-full whitespace-nowrap rounded-md border bg-white">
                                <div className="flex w-max space-x-4 p-4">
                                    {loading ? (
                                        <p className="text-generic">Loading...</p>
                                    ) : (
                                        group.items.map((item) => (
                                            <WardrobeItem key={item.id} item={item} />
                                        ))
                                    )}
                                </div>
                                <ScrollBar orientation="horizontal" />
                            </ScrollArea>
                        </div>
                    ))}
                </div>

                <div className="w-1/2 h-5/6 p-4">
                    <WardrobeOutfit />
                </div>
            </div>
        </DndProvider>
    );
};

export default WardrobePage;
