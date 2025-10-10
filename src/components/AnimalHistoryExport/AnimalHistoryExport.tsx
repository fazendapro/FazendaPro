import React from 'react';
import { Button, message } from 'antd';
import { ExportOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import jsPDF from 'jspdf';
import 'jspdf-autotable';
import { Animal } from '../../types/animal';
import { Sale } from '../../types/sale';
import { MilkCollection } from '../../types/milk-collection';
import { Reproduction } from '../../types/reproduction';

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
      const doc = new jsPDF();
      const pageWidth = doc.internal.pageSize.getWidth();
      const margin = 20;

      // Header
      doc.setFontSize(20);
      doc.setFont('helvetica', 'bold');
      doc.text(t('animalHistoryExport.title'), pageWidth / 2, 30, { align: 'center' });

      // Animal Info
      doc.setFontSize(14);
      doc.setFont('helvetica', 'bold');
      doc.text(t('animalHistoryExport.animalInfo'), margin, 50);

      doc.setFont('helvetica', 'normal');
      doc.setFontSize(12);
      let yPosition = 60;

      const animalInfo = [
        [t('animalHistoryExport.fields.name'), animal.animal_name],
        [t('animalHistoryExport.fields.earTag'), animal.ear_tag_number_local.toString()],
        [t('animalHistoryExport.fields.breed'), animal.breed],
        [t('animalHistoryExport.fields.type'), animal.type],
        [t('animalHistoryExport.fields.birthDate'), animal.birth_date ? new Date(animal.birth_date).toLocaleDateString('pt-BR') : t('common.notInformed')],
      ];

      animalInfo.forEach(([label, value]) => {
        doc.text(`${label}: ${value}`, margin, yPosition);
        yPosition += 8;
      });

      yPosition += 10;

      // Sales History
      if (sales.length > 0) {
        doc.setFont('helvetica', 'bold');
        doc.setFontSize(14);
        doc.text(t('animalHistoryExport.salesHistory'), margin, yPosition);
        yPosition += 10;

        const salesData = sales.map(sale => [
          new Date(sale.sale_date).toLocaleDateString('pt-BR'),
          sale.buyer_name,
          `R$ ${sale.price.toFixed(2)}`,
          sale.notes || '-'
        ]);

        (doc as any).autoTable({
          startY: yPosition,
          head: [[
            t('animalHistoryExport.salesTable.date'),
            t('animalHistoryExport.salesTable.buyer'),
            t('animalHistoryExport.salesTable.price'),
            t('animalHistoryExport.salesTable.notes')
          ]],
          body: salesData,
          styles: { fontSize: 10 },
          headStyles: { fillColor: [66, 139, 202] },
          margin: { left: margin, right: margin }
        });

        yPosition = (doc as any).lastAutoTable.finalY + 10;
      }

      // Milk Collections History
      if (milkCollections.length > 0) {
        doc.setFont('helvetica', 'bold');
        doc.setFontSize(14);
        doc.text(t('animalHistoryExport.milkHistory'), margin, yPosition);
        yPosition += 10;

        const milkData = milkCollections.map((collection: MilkCollection) => [
          new Date(collection.collection_date).toLocaleDateString('pt-BR'),
          `${collection.quantity}L`,
          collection.quality || '-',
          collection.notes || '-'
        ]);

        (doc as any).autoTable({
          startY: yPosition,
          head: [[
            t('animalHistoryExport.milkTable.date'),
            t('animalHistoryExport.milkTable.quantity'),
            t('animalHistoryExport.milkTable.quality'),
            t('animalHistoryExport.milkTable.notes')
          ]],
          body: milkData,
          styles: { fontSize: 10 },
          headStyles: { fillColor: [34, 139, 34] },
          margin: { left: margin, right: margin }
        });

        yPosition = (doc as any).lastAutoTable.finalY + 10;
      }

      // Reproduction History
      if (reproductions.length > 0) {
        doc.setFont('helvetica', 'bold');
        doc.setFontSize(14);
        doc.text(t('animalHistoryExport.reproductionHistory'), margin, yPosition);
        yPosition += 10;

        const reproductionData = reproductions.map((reproduction: Reproduction) => [
          new Date(reproduction.date).toLocaleDateString('pt-BR'),
          reproduction.phase,
          reproduction.notes || '-'
        ]);

        (doc as any).autoTable({
          startY: yPosition,
          head: [[
            t('animalHistoryExport.reproductionTable.date'),
            t('animalHistoryExport.reproductionTable.phase'),
            t('animalHistoryExport.reproductionTable.notes')
          ]],
          body: reproductionData,
          styles: { fontSize: 10 },
          headStyles: { fillColor: [255, 140, 0] },
          margin: { left: margin, right: margin }
        });
      }

      // Footer
      const pageCount = doc.getNumberOfPages();
      for (let i = 1; i <= pageCount; i++) {
        doc.setPage(i);
        doc.setFontSize(8);
        doc.setFont('helvetica', 'normal');
        doc.text(
          `${t('animalHistoryExport.footer.page')} ${i} ${t('animalHistoryExport.footer.of')} ${pageCount}`,
          pageWidth / 2,
          doc.internal.pageSize.getHeight() - 10,
          { align: 'center' }
        );
        doc.text(
          new Date().toLocaleDateString('pt-BR'),
          pageWidth - margin,
          doc.internal.pageSize.getHeight() - 10,
          { align: 'right' }
        );
      }

      // Save the PDF
      const fileName = `${animal.animal_name}_historico_${new Date().toISOString().split('T')[0]}.pdf`;
      doc.save(fileName);

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
