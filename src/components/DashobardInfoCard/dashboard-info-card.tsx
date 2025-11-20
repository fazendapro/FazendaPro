import { Card } from "antd";

export const DashboardInfoCard = ({title, value, icon, isLast}: {title: string, value: string, icon: React.ReactNode, isLast: boolean}) => {
  return (
      <Card style={{ textAlign: 'center', borderRadius: 0, border:'none', borderRight: isLast ? 'none' : '1px solid rgba(0,0,0,0.1)'}}>
          {icon}
          <p>{value}</p>
          <p>{title}</p>
        </Card>
  );
};