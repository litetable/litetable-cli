import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar.jsx";

export default function Layout({ children }) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <main className="p-4 w-full">
        <SidebarTrigger />
        {children}
      </main>
    </SidebarProvider>
  );
}
