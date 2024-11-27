import React, { useEffect, useState } from 'react';

function MyOrders() {
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    // Fetch orders (replace with API call)
    setOrders([
      { id: 1, status: 'Pending', details: 'Order 1' },
      { id: 2, status: 'In Transit', details: 'Order 2' }
    ]);
  }, []);

  return (
    <div>
      <h3>My Orders</h3>
      <ul>
        {orders.map(order => (
          <li key={order.id}>
            {order.details} - {order.status}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default MyOrders;