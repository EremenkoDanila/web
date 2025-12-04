import "./CartFloating.css";

interface Props {
  count: number;
  onOpen: () => void;
}

export const CartFloating = ({ count, onOpen }: Props) => {
  const handleClick = () => {
    if (count > 0) {
      onOpen();
    }
  };

  return (
    <div className="cart-floating" onClick={handleClick}>
      <img 
        src="http://127.0.0.1:9000/lb1/cart.png" 
        alt="cart" 
        className={`cart-icon ${count === 0 ? 'disabled' : ''}`} 
      />
      <div className={`cart-count ${count === 0 ? 'empty' : ''}`}>
        {count}
      </div>
    </div>
  );
};