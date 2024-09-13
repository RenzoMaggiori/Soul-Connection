import { FC } from 'react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";

interface DashboardCardProps {
    cardDescription: string;
    cardTitle: string;
    cardContent: string;
    progress: number;
    className?: string;
}

const DashboardCard: FC<DashboardCardProps> = ({ cardDescription, cardTitle, cardContent, progress, className }) => {
    return (
        <Card className={className}>
            <CardHeader className="pb-2">
                <CardDescription>{cardDescription}</CardDescription>
                <CardTitle className="text-4xl">{cardTitle}</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="text-xs text-muted-foreground">{cardContent}</div>
            </CardContent>
            <CardFooter>
                <Progress value={progress} aria-label={`${progress}% increase`} />
            </CardFooter>
        </Card>
    );
};

export default DashboardCard;
