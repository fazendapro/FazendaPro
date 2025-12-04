import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { DesktopTabs, MobileTabs, Tab } from '../../../../../components/tabs';
import { useIsMobile } from '../../../../../hooks/use-is-mobile';
import { Animals as AnimalsComponent } from "./animals";
import { MilkProduction } from "../milk-production/milk-production";
import { Reproduction } from "../reproduction/reproduction";
import { Vaccines } from "../vaccines/vaccines";

const Animals = () => {
  const isMobile = useIsMobile();
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState('animals-list');

  const tabsIndex = {
    'animals-list': 0,
    'milk-production': 1,
    'reproduction': 2,
    'vaccines': 3
  }

  const tabs: Tab[] = useMemo(() => [
    {
      title: t('animalTable.tabs.animalsList'),
      name: 'animals-list',
      component: (
        <AnimalsComponent />
      )
    },
    {
      title: t('animalTable.tabs.milkProduction'),
      name: 'milk-production',
      component: (
        <MilkProduction />
      )
    },
    {
      title: t('animalTable.tabs.reproduction'),
      name: 'reproduction',
      component: (
        <Reproduction />
      )
    },
    {
      title: t('animalTable.tabs.vaccines'),
      name: 'vaccines',
      component: (
        <Vaccines />
      )
    }
  ], [t]);

  const handleTabChange = (index: number) => {
    setActiveTab(tabs[index].name);
  };

  return (
    <div id="animals-container">
      {isMobile ? (
        <MobileTabs 
          tabs={tabs} 
          defaultTabIndex={tabsIndex[activeTab as keyof typeof tabsIndex] || 0}
          onChange={handleTabChange}
        />
      ) : (
        <DesktopTabs 
          tabs={tabs} 
          defaultTabIndex={tabsIndex[activeTab as keyof typeof tabsIndex] || 0}
          onChange={handleTabChange}
        />
      )}
    </div>
  );
};

export { Animals };