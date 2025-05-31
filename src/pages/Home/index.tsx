import { Layout } from "antd";
import { Footer } from "antd/es/layout/layout";


export const Home = () => {

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Footer style={{ textAlign: 'center' }}>
        <h1>
          FarmPro - Soluções em Agro ©{new Date().getFullYear()}
        </h1>
      </Footer>
    </Layout>
  );
};
