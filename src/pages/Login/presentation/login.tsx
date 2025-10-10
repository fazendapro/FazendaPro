import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { Button, Row, Col, Typography, Flex, Image } from 'antd';
import { useTranslation } from 'react-i18next';
import { Form } from '../../../components';
import { FieldType } from '../../../types/field-types';
import { loginSchema } from './login-schema';
import { toast } from 'react-toastify';
import { baseStyle } from './styles';
import { useAuth } from '../../../contexts/AuthContext';
import { useResponsive } from '../../../hooks';
import logo from '../../../assets/images/logo.png';

const { Title } = Typography;

interface LoginForm {
  email: string;
  password: string;
  remember?: boolean;
}

const Login = () => {
  const { login } = useAuth();
  const { t } = useTranslation();
  const { isMobile } = useResponsive();

  const methods = useForm<LoginForm>({
    resolver: yupResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
      remember: false,
    },
  });

  const onSubmit = async (data: LoginForm) => {
    try {
      await login(data.email, data.password);
    } catch (error) {
      toast.error(`Erro no login: ${error}`);
    }
  };

  const loginFields: FieldType[] = [
    {
      name: 'email',
      label: t('email'),
      type: 'text',
      placeholder: t('emailPlaceholder'),
      isRequired: true,
    },
    {
      name: 'password',
      label: t('password'),
      type: 'password',
      placeholder: t('passwordPlaceholder'),
      isRequired: true,
    },
    {
      name: 'forgotPassword',
      label: t('forgotPassword'),
      type: 'link',
    },
  ];

  return (
    <Flex 
      style={{ 
        padding: isMobile ? '20px' : '50px', 
        display: 'flex', 
        flexDirection: 'column', 
        justifyContent: 'center', 
        alignItems: 'center',
        minHeight: '100vh'
      }} 
      vertical 
      className="login-container"
    >
      <Row gutter={[16, 32]} style={{ width: '100%', maxWidth: '1200px' }}>
        <Col 
          className='logo-container' 
          xs={24} 
          sm={24} 
          md={12} 
          lg={12} 
          xl={12} 
          style={{ 
            textAlign: 'center',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center'
          }}
        >
          <div className="logo-placeholder">
            <div style={{ 
              fontSize: isMobile ? '120px' : '200px',
              maxWidth: '100%',
              height: 'auto'
            }}>
              <Image 
                src={logo} 
                preview={false} 
                alt="logo" 
                style={{ 
                  maxWidth: '100%', 
                  height: 'auto',
                  objectFit: 'contain'
                }}
              />
            </div>
          </div>
        </Col>

        <Col 
          xs={24} 
          sm={24} 
          md={12} 
          lg={12} 
          xl={12} 
          style={{
            ...baseStyle,
            padding: isMobile ? '20px' : '150px',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center'
          }}
        >
          <div className="form-container" style={{ width: '100%' }}>
            <Title level={3} style={{ textAlign: 'center', marginBottom: '24px' }}>
              {t('loginTitle')}
            </Title>
            <Form<LoginForm>
              onSubmit={onSubmit}
              fields={loginFields}
              methods={methods}
            >
              <Button
                className="login-button"
                type="primary"
                style={{ marginTop: '10px' }}
                htmlType="submit"
                block
                loading={methods.formState.isSubmitting}
                size="large"
              >
                {t('access')}
              </Button>
            </Form>
          </div>
        </Col>
      </Row>
    </Flex>
  );
};

export { Login };