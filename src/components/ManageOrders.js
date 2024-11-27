import React, { useEffect, useState } from 'react';
import UpdateOrderStatus from './UpdateOrderStatus';

function ManageOrders() {
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    // Fetch all orders (replace with API call)
    setOrders([
      { id: 1, status: 'Pending', details: 'Order 1' },
      { id: 2, status: 'In Transit', details: 'Order 2' }
    ]);
  }, []);

  const deleteOrder = async (orderId) => {
    try {
      await fetch(`http://localhost:5000/api/orders/${orderId}`, { method: 'DELETE' });
      setOrders(orders.filter(order => order.id !== orderId));
    } catch (error) {
      console.error('Error deleting order:', error);
    }
  };

  return (
    <div>
      <h3>Manage Orders</h3>
      <ul>
        {orders.map(order => (
          <li key={order.id}>
            {order.details} - {order.status}
            <UpdateOrderStatus orderId={order.id} />
            <button onClick={() => deleteOrder(order.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ManageOrders;    