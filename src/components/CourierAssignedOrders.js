import React, { useEffect, useState } from 'react';

function CourierAssignedOrders({ courierId }) {
  const [assignedOrders, setAssignedOrders] = useState([]);

  useEffect(() => {
    // Fetch assigned orders for a specific courier (replace with API call)
    setAssignedOrders([
      { id: 1, status: 'In Transit', details: 'Order 1' },
      { id: 2, status: 'Delivered', details: 'Order 2' }
    ]);
  }, [courierId]);

  return (
    <div>
      <h3>Assigned Orders</h3>
      <ul>
        {assignedOrders.map(order => (
          <li key={order.id}>
            {order.details} - {order.status}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default CourierAssignedOrders;