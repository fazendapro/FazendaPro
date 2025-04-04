import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { Button, Row, Col, Typography, Flex, Image } from 'antd';
import { useAuth } from '../../hooks/useAuth';
import { useTranslation } from 'react-i18next';
import { Form } from '../../components/lib/form';
import { FieldType } from '../../types/field-types';
import { loginSchema } from './login-schema';
import { toast } from 'react-toastify';
import logo from '../../assets/logo.png';
import { baseStyle } from './styles';

const { Title } = Typography;

interface LoginForm {
  email: string;
  password: string;
  remember?: boolean;
}

const Login = () => {
  const { login } = useAuth();
  const { t } = useTranslation();

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
      const success = await login(data.email, data.password);
      if (success) {
        toast.success(t('loginSuccess'));
      }
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
    <Flex vertical className="login-container">
      <Row gutter={16}>
        <Col className='logo-container' xs={24} md={12} style={{ textAlign: 'center' }}>
          <div className="logo-placeholder">
            <div style={{ fontSize: '200px' }}>
              <Image src={logo} preview={false} alt="logo" />
            </div>
          </div>
        </Col>

        <Col xs={24} md={12} style={baseStyle}>
          <div className="form-container">
            <Title level={3}>{t('loginTitle')}</Title>
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

export default Login;