import React, { useState } from 'react';

function CreateOrder() {
  const [pickup, setPickup] = useState('');
  const [dropoff, setDropoff] = useState('');
  const [details, setDetails] = useState('');
  const [time, setTime] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!pickup || !dropoff || !details) {
      alert('Please fill in all required fields');
      return;
    }
    console.log('Order created:', { pickup, dropoff, details, time });
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
