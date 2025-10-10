import { Grid } from "antd";
import { Sidebar } from "./sidebar";
import { MobileSidebar } from "./mobile-sidebar";

const { useBreakpoint } = Grid;

export const ResponsiveSidebar = () => {
  const screens = useBreakpoint();

  if (screens.xs) {
    return <MobileSidebar />;
  }

  return <Sidebar />;
};
