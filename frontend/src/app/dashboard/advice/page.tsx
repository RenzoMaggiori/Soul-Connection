"use client";

import { useQuery } from "@tanstack/react-query";
import { Tip } from "@/db/schemas";
import { getTips } from "@/db/tip";
import { Accordion, AccordionItem, AccordionTrigger, AccordionContent } from "@/components/ui/accordion";
import { useState } from "react";

export default function AdvicePage() {
  const { isLoading, isError, data } = useQuery({
    queryFn: async () => {
      const tips = await getTips();
      return { tips };
    },
    queryKey: ["TipsData"],
    gcTime: 1000 * 60, // 1 minute
  });

  const [selectedItem, setSelectedItem] = useState<string | null>(null);

  if (isError) return <div className="text-generic">Error loading tips...</div>;
  if (data && !data.tips) return <div className="text-generic">Parsing Error...</div>;

  const handleAccordionChange = (value: string | null) => {
    setSelectedItem(value);
  };

  return (
    <div className="mx-auto px-4 py-4 sm:px-6 lg:px-8">
      <h1 className="text-generic text-xl font-bold mb-6 md:text-2xl">Tips for Coaches</h1>
      <Accordion type="single" className="w-full" onValueChange={handleAccordionChange}>
        {data && data.tips?.map((tip: Tip) => (
          <AccordionItem
            key={tip.Id}
            value={`tip${tip.Id}`}
            className={`border rounded-md transition-colors bg-white shadow-lg`}>
            <AccordionTrigger className="text-generic font-semibold px-4 py-2">
              {tip.Title}
            </AccordionTrigger>
            <AccordionContent className="text-generic leading-relaxed px-4 py-2">
              {tip.Tip}
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  );
}

// "use client";

// import { useQuery } from "@tanstack/react-query";
// import { Tip, tipSchema } from "@/db/schemas";
// import { getTips } from "@/db/tip";
// import { useEffect, useState } from "react";

// export default function AdvicePage() {
//   const { isLoading, isError, data } = useQuery({
//     queryFn: async () => {
//       const tips = await getTips();
//       return { tips };
//     },
//     queryKey: ["TipsData"],
//     gcTime: 1000 * 60, // 1 minute
//   });

//   if (isError) return <div className="text-generic">Error loading tips...</div>;
//   if (data && !data.tips) return <div className="text-generic">Parsing Error...</div>;

//   return (
//     <div className="mx-auto px-4 py-10 sm:px-6 lg:px-8">
//       <div className="grid max-w-7xl mx-auto grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
//         {data && data.tips?.map((tip: Tip) => (
//           <div
//             key={tip.Id}
//             className="bg-white shadow-lg rounded-lg p-6 hover:shadow-xl transition-shadow duration-300 ease-in-out"
//           >
//             <h3 className="text-2xl font-semibold text-generic mb-3">{tip.Title}</h3>
//             <p className="text-generic text-lg leading-relaxed">{tip.Tip}</p>
//           </div>
//         ))}
//       </div>
//     </div>
//   );
// }

