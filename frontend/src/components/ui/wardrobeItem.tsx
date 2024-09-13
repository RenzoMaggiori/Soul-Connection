import React from "react";
import Image from "next/image";
import { useDrag } from "react-dnd";

interface WardrobeItemProps {
    item: {
        id: number;
        image: string;
        type: string;
    };
}

const WardrobeItem: React.FC<WardrobeItemProps> = ({ item }) => {
    const [, dragRef] = useDrag(() => ({
        type: "wardrobeItem",
        item,
    }));

    return (
        <div ref={dragRef} className="border rounded-md p-4 cursor-pointer w-fit h-fit inline-block bg-secondary">
            <Image
                src={item.image}
                alt={`Clothing id: ${item.id}\nType: ${item.type}`}
                width={100}
                height={100}
                className="rounded"
            />
        </div>
    );
};

export default WardrobeItem;
