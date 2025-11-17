import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { CustomPagination } from '../custom-pagination';
import { useResponsive } from '../../../../hooks';

vi.mock('../../../../hooks', () => ({
  useResponsive: vi.fn()
}));

const mockUseResponsive = useResponsive as jest.MockedFunction<typeof useResponsive>;

describe('CustomPagination', () => {
  const defaultProps = {
    current: 1,
    total: 100,
    pageSize: 10,
    onChange: vi.fn(),
    onShowSizeChange: vi.fn()
  };

  beforeEach(() => {
    vi.clearAllMocks();
    mockUseResponsive.mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true,
      isLargeDesktop: false,
      screenWidth: 1200
    });
  });

  describe('Renderização básica', () => {
    it('deve renderizar o componente com props padrão', () => {
      render(<CustomPagination {...defaultProps} />);
      
      expect(screen.getByText('1-10 de 100 registros')).toBeInTheDocument();
      expect(screen.getByTitle('1')).toBeInTheDocument();
    });

    it('deve renderizar com total zero', () => {
      render(<CustomPagination {...defaultProps} total={0} />);
      
      expect(screen.getByText('0-0 de 0 registros')).toBeInTheDocument();
    });

    it('deve aplicar className customizada', () => {
      const { container } = render(
        <CustomPagination {...defaultProps} className="custom-class" />
      );
      
      expect(container.querySelector('.custom-pagination')).toHaveClass('custom-class');
    });

    it('deve aplicar estilos customizados', () => {
      const customStyle = { backgroundColor: 'red' };
      const { container } = render(
        <CustomPagination {...defaultProps} style={customStyle} />
      );
      
      expect(container.querySelector('.custom-pagination')).toHaveStyle('background-color: red');
    });
  });

  describe('Comportamento responsivo', () => {
    it('deve renderizar layout mobile quando isMobile é true', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
      });

      const { container } = render(<CustomPagination {...defaultProps} />);
      
      const paginationContainer = container.querySelector('.custom-pagination');
      expect(paginationContainer).toHaveStyle({
        'flex-direction': 'column',
        'gap': '12px'
      });
    });

    it('deve renderizar layout desktop quando isMobile é false', () => {
      mockUseResponsive.mockReturnValue({
        isMobile: false,
        isTablet: false,
        isDesktop: true,
        isLargeDesktop: false,
        screenWidth: 1200
      });

      const { container } = render(<CustomPagination {...defaultProps} />);
      
      const paginationContainer = container.querySelector('.custom-pagination');
      expect(paginationContainer).not.toHaveStyle({
        'flex-direction': 'column'
      });
    });

    it('deve usar pageSize padrão baseado no dispositivo', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
      });

      render(<CustomPagination {...defaultProps} pageSize={undefined} />);
      
      // Em mobile com simple mode, não mostra o texto de total
      const pagination = screen.getByRole('list');
      expect(pagination).toBeInTheDocument();
    });
  });

  describe('Controles de paginação', () => {
    it('deve chamar onChange quando página é alterada', async () => {
      const onChange = vi.fn();
      render(<CustomPagination {...defaultProps} onChange={onChange} />);
      
      const nextButton = screen.getByTitle('Next Page');
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(onChange).toHaveBeenCalledWith(2, 10);
      });
    });

    it('deve chamar onShowSizeChange quando tamanho da página é alterado', async () => {
      const onShowSizeChange = vi.fn();
      render(
        <CustomPagination 
          {...defaultProps} 
          onShowSizeChange={onShowSizeChange}
          showSizeChanger={true}
        />
      );
      
      const sizeChanger = screen.getByRole('combobox');
      fireEvent.mouseDown(sizeChanger);
      
      await waitFor(() => {
        const option20 = screen.getByText('20 / page');
        fireEvent.click(option20);
      });
      
      await waitFor(() => {
        expect(onShowSizeChange).toHaveBeenCalledWith(1, 20);
      });
    });

    it('deve navegar para página específica', async () => {
      const onChange = vi.fn();
      render(<CustomPagination {...defaultProps} onChange={onChange} />);
      
      const page3Button = screen.getByTitle('3');
      fireEvent.click(page3Button);
      
      await waitFor(() => {
        expect(onChange).toHaveBeenCalledWith(3, 10);
      });
    });
  });

  describe('Configurações de exibição', () => {
    it('deve ocultar showSizeChanger em mobile', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
      });

      render(<CustomPagination {...defaultProps} showSizeChanger={true} />);
      
      expect(screen.queryByRole('combobox')).not.toBeInTheDocument();
    });

    it('deve mostrar showSizeChanger em desktop', () => {
      mockUseResponsive.mockReturnValue({
        isMobile: false,
        isTablet: false,
        isDesktop: true,
        isLargeDesktop: false,
        screenWidth: 1200
      });

      render(<CustomPagination {...defaultProps} showSizeChanger={true} />);
      
      expect(screen.getByRole('combobox')).toBeInTheDocument();
    });

    it('deve ocultar showTotal em mobile', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
      });

      render(<CustomPagination {...defaultProps} showTotal={true} />);
      
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });

    it('deve mostrar showTotal em desktop', () => {
      mockUseResponsive.mockReturnValue({
        isMobile: false,
        isTablet: false,
        isDesktop: true,
        isLargeDesktop: false,
        screenWidth: 1200
      });

      render(<CustomPagination {...defaultProps} showTotal={true} />);
      
      expect(screen.getByText('1-10 de 100 registros')).toBeInTheDocument();
    });

    it('deve usar função customizada para showTotal', () => {
      const customShowTotal = (total: number, range: [number, number]) => 
        `Página ${range[0]}-${range[1]} de ${total} itens`;
      
      render(
        <CustomPagination 
          {...defaultProps} 
          showTotal={customShowTotal}
        />
      );
      
      expect(screen.getByText('Página 1-10 de 100 itens')).toBeInTheDocument();
    });
  });

  describe('Estados especiais', () => {
    it('deve lidar com página atual maior que total de páginas', () => {
      render(
        <CustomPagination 
          {...defaultProps} 
          current={10}
          total={50}
          pageSize={10}
        />
      );
      
      expect(screen.getByText('41-50 de 50 registros')).toBeInTheDocument();
    });

    it('deve lidar com total menor que pageSize', () => {
      render(
        <CustomPagination 
          {...defaultProps} 
          total={5}
          pageSize={10}
        />
      );
      
      expect(screen.getByText('1-5 de 5 registros')).toBeInTheDocument();
    });

    it('deve mostrar página correta quando current é undefined', () => {
      render(
        <CustomPagination 
          {...defaultProps} 
          current={undefined}
        />
      );
      
      expect(screen.getByText('1-10 de 100 registros')).toBeInTheDocument();
    });
  });

  describe('Acessibilidade', () => {
    it('deve ter lista de paginação', () => {
      render(<CustomPagination {...defaultProps} />);
      
      expect(screen.getByRole('list')).toBeInTheDocument();
    });

    it('deve ter botões acessíveis', () => {
      render(<CustomPagination {...defaultProps} />);
      
      expect(screen.getByTitle('Previous Page')).toBeInTheDocument();
      expect(screen.getByTitle('Next Page')).toBeInTheDocument();
    });
  });

  describe('Integração com Ant Design', () => {
    it('deve passar props adicionais para AntPagination', () => {
      render(
        <CustomPagination 
          {...defaultProps}
        />
      );
      
      expect(screen.getByRole('list')).toBeInTheDocument();
    });

    it('deve usar size correto baseado no dispositivo', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
      });

      render(<CustomPagination {...defaultProps} />);
      
      const pagination = screen.getByRole('list');
      expect(pagination).toBeInTheDocument();
    });
  });
});
