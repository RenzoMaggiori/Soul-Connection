"use client";

import React, { useEffect, useRef } from "react";

interface MapsComponentProps {
    lat: number;
    lng: number;
}

declare global {
    interface Window {
        initMap: () => void;
    }
}

const MapsComponent: React.FC<MapsComponentProps> = ({ lat, lng }) => {
    const mapRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        const loadGoogleMapsScript = () => {
            if (!document.getElementById("google-maps-script")) {
                const script = document.createElement("script");
                script.id = "google-maps-script";
                script.src = `https://maps.googleapis.com/maps/api/js?key=YOUR_GOOGLE_MAPS_API_KEY&callback=initMap`;
                script.async = true;
                script.defer = true;
                document.head.appendChild(script);
            }
        };

        window.initMap = () => {
            if (mapRef.current) {
                new google.maps.Map(mapRef.current, {
                    center: { lat, lng },
                    zoom: 15,
                });
            }
        };

        loadGoogleMapsScript();

        return () => {
            const existingScript = document.getElementById("google-maps-script");
            if (existingScript) existingScript.remove();
        };
    }, [lat, lng]);

    return <div ref={mapRef} style={{ width: "100%", height: "500px" }} />;
};

export default MapsComponent;
