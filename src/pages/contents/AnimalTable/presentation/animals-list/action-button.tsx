import React from 'react';

interface ActionButtonProps {
  animalId: string;
}

export const ActionButton: React.FC<ActionButtonProps> = ({ animalId }) => {
  const handleClick = () => {
    window.location.href = `/animal/${animalId}`;
  };

  return (
    <button 
      onClick={handleClick}
      style={{
        background: 'none',
        border: 'none',
        color: '#1890ff',
        cursor: 'pointer',
        textDecoration: 'underline'
      }}
    >
      Ver Ficha
    </button>
  );
};
