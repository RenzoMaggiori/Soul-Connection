"use client";

import React from "react";
import { MapContainer, TileLayer, Marker, Popup, useMap } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import { useQuery } from "@tanstack/react-query";

import { Event } from "@/db/schemas";
import { getEvents } from "@/db/event";
import { NewEventDialog } from "@/app/dashboard/events/newEventDialog";
import { barcelonaCoordinates } from "@/app/dashboard/events/utils";
import Calendar from "@/components/calendar/calendarComponent";
import { CalendarHeader } from "./calendar-header";
import { CalendarProvider } from "@/components/calendar/utils/calendar-context";

const customMarkerIcon = new L.Icon({
  iconUrl: "https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png",
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowUrl: "https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png",
  shadowSize: [41, 41],
});

// Component to handle the map view
function MapView({ selectedEvent }: { selectedEvent: Event | null }) {
  const map = useMap();

  React.useEffect(() => {
    if (selectedEvent) {
      map.flyTo(
        [
          parseFloat(selectedEvent.Location_X),
          parseFloat(selectedEvent.Location_Y),
        ],
        15,
      );
    }
  }, [selectedEvent, map]);

  return null;
}

export default function EventsPage() {
  const [selectedEvent, setSelectedEvent] = React.useState<Event | null>(null);
  const markersRef = React.useRef<Map<number, L.Marker>>(new Map());

  const { isLoading, isError, data } = useQuery({
    queryFn: async () => {
      const events = await getEvents();
      if (!events) {
        throw new Error("Something went wrong");
      }
      return { events };
    },
    queryKey: ["EventsData"],
    gcTime: 1000 * 60, // 1 minute
  });

  React.useEffect(() => {
    if (selectedEvent && markersRef.current.has(selectedEvent.Id)) {
      const marker = markersRef.current.get(selectedEvent.Id);
      if (marker) {
        marker.openPopup();
      }
    }
  }, [selectedEvent]);

  if (isError) return <div>Error...</div>;
  if (!data || isLoading) return <div>Loading...</div>;
  if (data.events.length === 0) return <div>No events found...</div>;

  return (
    <div className="flex h-full flex-col gap-4 p-2">
      {/* Calendar */}
      <div className="h-[45vh]">
        <CalendarProvider>
          <Calendar
            events={data.events}
            eventOnClick={setSelectedEvent}
            Header={CalendarHeader}
          />
        </CalendarProvider>
      </div>
      {/* Map */}
      <div className="h-[39vh]">
        <MapContainer
          center={[barcelonaCoordinates.x, barcelonaCoordinates.y]}
          zoom={10}
          className="z-10 h-full w-full"
        >
          <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
          <MapView selectedEvent={selectedEvent} />
          {data.events.map((event) => (
            <Marker
              key={event.Id}
              position={[
                parseFloat(event.Location_X),
                parseFloat(event.Location_Y),
              ]}
              icon={customMarkerIcon}
              ref={(marker) => {
                if (marker) {
                  markersRef.current.set(event.Id, marker);
                }
              }}
            >
              <Popup>
                <div className="flex items-start justify-between">
                  <div className="text-generic">
                    <strong className="text-base">{event.Name}</strong>
                    <br />
                    Max Participants: {event.Max_Participants}
                  </div>
                  <span className="ml-4 mt-1 whitespace-nowrap text-generic">
                    {event.Date}
                  </span>
                </div>
              </Popup>
            </Marker>
          ))}
        </MapContainer>
      </div>
    </div>
  );
}
