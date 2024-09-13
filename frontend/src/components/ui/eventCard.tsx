import React from "react";
import { MapPin, Users } from "lucide-react";

interface EventCardProps {
    name: string;
    date: string;
    duration?: number;
    maxParticipants: number;
    locationX?: number;
    locationY?: number;
    type?: string;
    employeeId?: number;
    locationName: string;
}

const EventCard: React.FC<EventCardProps> = ({ name, date, locationName, maxParticipants }) => {
    return (
        <div className="border-b p-2 flex justify-between items-center">
            <div>
                <h2 className="text-xl font-bold">{name}</h2>
                <div className="flex items-center">
                    <MapPin size={16}  className="mr-2" />
                    <p>{locationName}</p>
                </div>
                <div className="flex items-center">
                    <Users size={16} className="mr-2" />
                    <p>Max participants: {maxParticipants}</p>
                </div>
            </div>
            <p>{date}</p>
        </div>
    );
};

export default EventCard;
