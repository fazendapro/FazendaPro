import { RemoteGetMonthlySalesAndPurchases } from '../../data/usecases/remote-get-monthly-sales-and-purchases';

export const GetMonthlySalesAndPurchasesFactory = () => {
  return new RemoteGetMonthlySalesAndPurchases();
};


