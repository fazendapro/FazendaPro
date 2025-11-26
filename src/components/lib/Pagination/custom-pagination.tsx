import React from 'react';
import { Pagination as AntPagination } from 'antd';
import { useResponsive } from '../../../hooks';

interface CustomPaginationProps {
  current?: number;
  total?: number;
  pageSize?: number;
  showSizeChanger?: boolean;
  showTotal?: boolean | ((total: number, range: [number, number]) => string);
  onChange?: (page: number, pageSize: number) => void;
  onShowSizeChange?: (current: number, size: number) => void;
  className?: string;
  style?: React.CSSProperties;
}

export const CustomPagination: React.FC<CustomPaginationProps> = ({
  current,
  total = 0,
  pageSize = 10,
  showSizeChanger = true,
  showTotal = true,
  onChange,
  onShowSizeChange,
  className = '',
  style = {},
  ...props
}) => {
  const { isMobile, isTablet } = useResponsive();

  const defaultPageSize = 10;
  const finalPageSize = pageSize || defaultPageSize;

  const defaultShowTotal = (total: number, range: [number, number]) => 
    `${range[0]}-${range[1]} de ${total} registros`;

  const paginationStyle: React.CSSProperties = {
    marginTop: '24px',
    padding: '16px 0',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#fafafa',
    borderRadius: '8px',
    border: '1px solid #f0f0f0',
    ...style
  };

  const mobilePaginationStyle: React.CSSProperties = {
    ...paginationStyle,
    flexDirection: 'column',
    gap: '12px',
    padding: '12px',
    marginTop: '16px'
  };

  return (
    <div 
      className={`custom-pagination ${className}`}
      style={isMobile ? mobilePaginationStyle : paginationStyle}
    >
      <AntPagination
        current={current}
        total={total}
        pageSize={finalPageSize}
        showSizeChanger={!isMobile && showSizeChanger}
        showTotal={!isMobile && showTotal ? (showTotal === true ? defaultShowTotal : showTotal) : undefined}
        onChange={onChange}
        onShowSizeChange={onShowSizeChange}
        size={isMobile ? 'small' : 'default'}
        showLessItems={isMobile}
        simple={isMobile}
        style={{
          fontSize: isMobile ? '12px' : '14px'
        }}
        {...props}
      />
    </div>
  );
};
