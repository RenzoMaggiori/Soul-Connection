import * as React from "react";
import Image from "next/image";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Dialog, DialogContent, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useLoadScript } from "@react-google-maps/api";
import Maps from "@/components/ui/maps";

interface Event {
    name: string;
    date: string;
    location: string;
    image: string;
    attendees: number;
    locationX: number;
    locationY: number;
}

interface ScrollAreaEventsProps {
    title: string;
    events: Event[];
}

const ScrollAreaEvents: React.FC<ScrollAreaEventsProps> = ({ title, events }) => {
    const { isLoaded } = useLoadScript({
        googleMapsApiKey: "API_KEY",
    });

    const [selectedEvent, setSelectedEvent] = React.useState<Event | null>(null);

    if (!isLoaded) return <div>Loading...</div>;

    return (
        <div className="mb-6">
            <h2 className="text-lg font-semibold mb-2">{title}</h2>
            <ScrollArea className="w-full whitespace-nowrap rounded-md border">
                <div className="flex w-max space-x-4 p-4">
                    {events.map((event) => (
                        <figure key={event.name} className="shrink-0">
                            <Dialog>
                                <DialogTrigger asChild>
                                    <div className="overflow-hidden rounded-md cursor-pointer" onClick={() => setSelectedEvent(event)}>
                                        <Image
                                            src={event.image}
                                            alt={event.name}
                                            className="aspect-[3/4] h-fit w-fit object-cover"
                                            width="100"
                                            height="200"
                                        />
                                    </div>
                                </DialogTrigger>
                                <DialogContent className="sm:max-w-md">
                                    <DialogTitle>{event.name}</DialogTitle>
                                    <p className="text-sm text-muted-foreground">{event.date}</p>
                                    <p className="text-sm text-muted-foreground">Attendees: {event.attendees}</p>
                                    <p className="text-sm text-muted-foreground">{event.location}</p>
                                    {selectedEvent && (
                                        <Maps lat={event.locationY} lng={event.locationX} />
                                    )}
                                </DialogContent>
                            </Dialog>
                            <figcaption className="pt-2 text-xs text-muted-foreground">
                                <div className="font-semibold text-foreground">{event.name}</div>
                                <div>{event.date}</div>
                            </figcaption>
                        </figure>
                    ))}
                </div>
                <ScrollBar orientation="horizontal" />
            </ScrollArea>
        </div>
    );
};

export default ScrollAreaEvents;
