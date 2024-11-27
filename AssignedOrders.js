import React, { useState } from 'react';

function AssignedOrders() {
  const [orders, setOrders] = useState([
    { id: 1, details: 'Order 1', status: 'Assigned' },
    { id: 2, details: 'Order 2', status: 'Assigned' }
  ]);

  const updateStatus = (id, status) => {
    setOrders(orders.map(order => (order.id === id ? { ...order, status } : order)));
  };

  return (
    <div>
      <h3>Assigned Orders</h3>
      <ul>
        {orders.map(order => (
          <li key={order.id}>
            {order.details} - {order.status}
            <button onClick={() => updateStatus(order.id, 'Accepted')}>Accept</button>
            <button onClick={() => updateStatus(order.id, 'Declined')}>Decline</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default AssignedOrders;