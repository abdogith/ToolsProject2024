import React, { useState } from 'react';

function CreateOrder() {
  const [pickup, setPickup] = useState('');
  const [dropoff, setDropoff] = useState('');
  const [details, setDetails] = useState('');
  const [time, setTime] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const orderData = { pickup, dropoff, details, time };

    try {
      const response = await fetch('http://localhost:5000/api/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(orderData),
      });

      if (response.ok) {
        alert('Order created successfully');
      } else {
        alert('Failed to create order');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input type="text" placeholder="Pickup Location" value={pickup} onChange={(e) => setPickup(e.target.value)} required />
      <input type="text" placeholder="Drop-off Location" value={dropoff} onChange={(e) => setDropoff(e.target.value)} required />
      <input type="text" placeholder="Package Details" value={details} onChange={(e) => setDetails(e.target.value)} required />
      <input type="datetime-local" value={time} onChange={(e) => setTime(e.target.value)} />
      <button type="submit">Create Order</button>
    </form>
  );
}

export default CreateOrder;