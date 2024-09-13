import { DataTable } from "@/components/ui/data-table";
import { Task } from "@/app/dashboard/task_management/task";
import { TaskTableFooter, TaskTableHeader, taskColumns } from "@/app/dashboard/task_management/columns";

export default function TaskManagementPage() {
    const data: Task[] = [
        {
            title: "First task",
            description: "Finish the web",
            status: "doing",
            dueDate: new Date("1-1-2025")
        },
        {
            title: "First task",
            description: "Finish the web",
            status: "doing",
            dueDate: new Date("1-1-2025")
        },
        {
            title: "First task",
            description: "Finish the web",
            status: "doing",
            dueDate: new Date("1-1-2025")
        },
        {
            title: "First task",
            description: "Finish the web",
            status: "doing",
            dueDate: new Date("1-1-2025")
        },
        {
            title: "First task",
            description: "Finish the web",
            status: "doing",
            dueDate: new Date("1-1-2025")
        },
        {
            title: "First task",
            description: "Finish the web",
            status: "doing",
            dueDate: new Date("1-1-2025")
        }
    ]

    return (
        <div className="w-full h-full max-h-full pb-4">
            <div>
            <DataTable data={data.concat(data.concat(data))} columns={taskColumns} header={TaskTableHeader} footer={TaskTableFooter} pageSize={8}/>
            </div>
        </div>
    );
}
