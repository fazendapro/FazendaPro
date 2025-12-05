import React from 'react';
import { Button, message } from 'antd';
import { ExportOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Animal } from '../../pages/contents/AnimalTable/types/type';
import { Sale } from '../../types/sale';
import { MilkCollection } from '../../types/milk-collection';
import { Reproduction } from '../../types/reproduction';
import { generateAnimalHistoryPDF } from './pdfGenerator';

interface AnimalHistoryExportProps {
  animal: Animal;
  sales: Sale[];
  milkCollections?: MilkCollection[];
  reproductions?: Reproduction[];
}

export const AnimalHistoryExport: React.FC<AnimalHistoryExportProps> = ({
  animal,
  sales,
  milkCollections = [],
  reproductions = []
}) => {
  const { t } = useTranslation();

  const generatePDF = () => {
    try {
      generateAnimalHistoryPDF({
        animal,
        sales,
        milkCollections,
        reproductions,
        animalImage: animal.photo,
        translations: {
          animalHistoryExport: {
            title: t('animalHistoryExport.title'),
            animalInfo: t('animalHistoryExport.animalInfo'),
            fields: {
              name: t('animalHistoryExport.fields.name'),
              earTag: t('animalHistoryExport.fields.earTag'),
              breed: t('animalHistoryExport.fields.breed'),
              type: t('animalHistoryExport.fields.type'),
              birthDate: t('animalHistoryExport.fields.birthDate'),
              sex: t('animalHistoryExport.fields.sex'),
              status: t('animalHistoryExport.fields.status'),
              confinement: t('animalHistoryExport.fields.confinement'),
              fertilization: t('animalHistoryExport.fields.fertilization'),
              castrated: t('animalHistoryExport.fields.castrated'),
              currentWeight: t('animalHistoryExport.fields.currentWeight'),
              idealWeight: t('animalHistoryExport.fields.idealWeight'),
              milkProduction: t('animalHistoryExport.fields.milkProduction')
            },
            sex: {
              male: t('animalHistoryExport.sex.male'),
              female: t('animalHistoryExport.sex.female')
            },
            status: {
              active: t('animalHistoryExport.status.active'),
              inactive: t('animalHistoryExport.status.inactive'),
              sold: t('animalHistoryExport.status.sold'),
              deceased: t('animalHistoryExport.status.deceased')
            },
            statistics: t('animalHistoryExport.statistics'),
            stats: {
              totalSales: t('animalHistoryExport.stats.totalSales'),
              totalSalesValue: t('animalHistoryExport.stats.totalSalesValue'),
              totalMilkCollections: t('animalHistoryExport.stats.totalMilkCollections'),
              totalMilkQuantity: t('animalHistoryExport.stats.totalMilkQuantity'),
              totalReproductions: t('animalHistoryExport.stats.totalReproductions')
            },
            salesHistory: t('animalHistoryExport.salesHistory'),
            salesTable: {
              date: t('animalHistoryExport.salesTable.date'),
              buyer: t('animalHistoryExport.salesTable.buyer'),
              price: t('animalHistoryExport.salesTable.price'),
              notes: t('animalHistoryExport.salesTable.notes')
            },
            milkHistory: t('animalHistoryExport.milkHistory'),
            milkTable: {
              date: t('animalHistoryExport.milkTable.date'),
              quantity: t('animalHistoryExport.milkTable.quantity'),
              quality: t('animalHistoryExport.milkTable.quality'),
              notes: t('animalHistoryExport.milkTable.notes')
            },
            reproductionHistory: t('animalHistoryExport.reproductionHistory'),
            reproductionTable: {
              date: t('animalHistoryExport.reproductionTable.date'),
              phase: t('animalHistoryExport.reproductionTable.phase'),
              notes: t('animalHistoryExport.reproductionTable.notes')
            },
            footer: {
              page: t('animalHistoryExport.footer.page'),
              of: t('animalHistoryExport.footer.of')
            },
            units: {
              kg: t('animalHistoryExport.units.kg'),
              liters: t('animalHistoryExport.units.liters'),
              litersPerDay: t('animalHistoryExport.units.litersPerDay'),
              currency: t('animalHistoryExport.units.currency')
            },
            dateFormat: t('animalHistoryExport.dateFormat'),
            fileName: t('animalHistoryExport.fileName')
          },
          common: {
            notInformed: t('common.notInformed'),
            yes: t('common.yes'),
            no: t('common.no')
          }
        }
      });

      message.success(t('animalHistoryExport.success'));
    } catch (error) {
      console.error('Error generating PDF:', error);
      message.error(t('animalHistoryExport.error'));
    }
  };

  return (
    <Button
      icon={<ExportOutlined />}
      onClick={generatePDF}
      style={{
        backgroundColor: '#f0f0f0',
        borderColor: '#d9d9d9',
        color: '#262626',
        borderRadius: '6px'
      }}
    >
      {t('animalDetail.exportHistory')}
    </Button>
  );
};