"use client";

import Link from "next/link";

import { useShallow } from "zustand/react/shallow";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { APP_CONFIG } from "@/config/app-config";
import { useCurrentUser } from "@/hooks/use-current-user";
import { sidebarItems } from "@/navigation/sidebar/sidebar-items";
import { usePreferencesStore } from "@/stores/preferences/preferences-provider";

import { NavMain } from "./nav-main";
import { NavUser } from "./nav-user";
import { SidebarSupportCard } from "./sidebar-support-card";


export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const currentUser = useCurrentUser();
  const { sidebarVariant, sidebarCollapsible, isSynced } = usePreferencesStore(
    useShallow((s) => ({
      sidebarVariant: s.sidebarVariant,
      sidebarCollapsible: s.sidebarCollapsible,
      isSynced: s.isSynced,
    })),
  );

  const variant = isSynced ? sidebarVariant : props.variant;
  const collapsible = isSynced ? sidebarCollapsible : props.collapsible;

  return (
    <Sidebar {...props} variant={variant} collapsible={collapsible}>
      <SidebarHeader className="border-b border-sidebar-border/70 px-3 py-3">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild className="h-11 rounded-lg px-2 hover:bg-sidebar-accent/65">
              <Link prefetch={false} href="/dashboard">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src="/logo.jpg" alt="ARTEX" width={34} height={34} className="size-[34px] shrink-0 rounded-lg object-cover shadow-sm" />
                <span className="flex min-w-0 flex-col leading-tight">
                  <strong className="text-sm font-semibold tracking-wide">{APP_CONFIG.name}</strong>
                  <span className="text-[10px] font-normal text-sidebar-foreground/60">智能工作台</span>
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent className="py-2">
        <NavMain items={sidebarItems} />
        {/* <NavDocuments items={data.documents} /> */}
        {/* <NavSecondary items={data.navSecondary} className="mt-auto" /> */}
      </SidebarContent>
      <SidebarFooter className="border-t border-sidebar-border/70 p-3">
        <NavUser user={currentUser} />
      </SidebarFooter>
    </Sidebar>
  );
}
