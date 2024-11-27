import React, { useState } from 'react';

function UpdateOrderStatus({ orderId }) {
  const [status, setStatus] = useState('');

  useEffect(() => {
    // Fetch initial status (if required)
  }, [orderId]);

  const handleUpdateStatus = async () => {
    try {
      const response = await fetch(`http://localhost:5000/api/orders/${orderId}/status`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status }),
      });

      if (response.ok) {
        alert('Order status updated successfully');
      } else {
        alert('Failed to update status');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div>
      <select value={status} onChange={(e) => setStatus(e.target.value)}>
        <option value="">Select Status</option>
        <option value="Picked Up">Picked Up</option>
        <option value="In Transit">In Transit</option>
        <option value="Delivered">Delivered</option>
      </select>
      <button onClick={handleUpdateStatus}>Update Status</button>
    </div>
  );
}

export default UpdateOrderStatus;