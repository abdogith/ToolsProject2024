import React from 'react';
import { useParams } from 'react-router-dom';

function OrderDetails() {
  const { orderId } = useParams();
  const order = { id: orderId, pickup: 'Location A', dropoff: 'Location B', status: 'Pending' };

  const cancelOrder = () => {
    if (order.status === 'Pending') {
      alert('Order cancelled');
    }
  };

  return (
    <div>
      <h3>Order Details</h3>
      <p>Pickup: {order.pickup}</p>
      <p>Drop-off: {order.dropoff}</p>
      <p>Status: {order.status}</p>
      {order.status === 'Pending' && <button onClick={cancelOrder}>Cancel Order</button>}
    </div>
  );
}

export default OrderDetails;
