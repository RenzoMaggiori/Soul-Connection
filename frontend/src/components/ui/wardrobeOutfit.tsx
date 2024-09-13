import React, { useState } from "react";
import { useDrop } from "react-dnd";
import { Button } from "./button";
import { Info } from "lucide-react";

const WardrobeOutfit = () => {
  const [outfit, setOutfit] = useState<{ id: number; image: string }[]>([]);

  const [, dropRef] = useDrop(() => ({
    accept: "wardrobeItem",
    drop: (item: { id: number; image: string }) => {
      setOutfit((prevOutfit) => [...prevOutfit, item]);
    },
  }));

  return (
    <>
      <div className="flex flex-row items-center">
        <h2 className="text-2xl font-semibold text-generic">Outfit Builder</h2>
        <Button
          variant="generic"
          className="ml-auto block"
          onClick={() => setOutfit([])}
        >
          Clear Outfit
        </Button>
      </div>
      <div className="flex items-center text-generic mb-2">
          <Info className="mr-2" size={18} />
          <p>Drag and drop items from the wardrobe to build the outfit</p>
      </div>
      <div ref={dropRef} className="mt-4 h-full rounded-md border p-4 bg-white">
        <div className="mt-4 flex flex-wrap gap-4">
          {outfit.map((item, index) => (
            <div key={index} className="rounded-md border p-2 bg-secondary">
              <img
                src={item.image}
                alt={`Outfit item ${item.id}`}
                width={100}
                height={100}
              />
            </div>
          ))}
        </div>
      </div>
    </>
  );
};

export default WardrobeOutfit;
