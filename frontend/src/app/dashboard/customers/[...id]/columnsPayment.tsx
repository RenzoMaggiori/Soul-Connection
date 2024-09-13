import { ColumnDef } from "@tanstack/react-table";
import { Payment } from "./payment"

export const paymentColumns: ColumnDef<Payment>[] = [
    {
        accessorKey: "Date",
        header: "Date",
        cell: ({ row }) => <div>{row.original.Date}</div>,
    },
    {
        accessorKey: "PaymentMethod",
        header: "Payment Method",
        cell: ({ row }) => <div>{row.original.PaymentMethod}</div>,
    },
    {
        accessorKey: "Amount",
        header: "Amount",
        cell: ({ row }) => <div>{row.original.Amount.toFixed(2)} €</div>,
    },
    {
        accessorKey: "Comment",
        header: "Comment",
        cell: ({ row }) => <div>{row.original.Comment}</div>,
    },
];
