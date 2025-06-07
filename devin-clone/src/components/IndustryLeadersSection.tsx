import React from 'react';

const IndustryLeadersSection = () => {
  return (
    <section className="py-12 md:py-24 bg-gray-800 text-white">
      <div className="container mx-auto px-6 text-center">
        <h2 className="text-3xl md:text-4xl font-bold mb-8">
          Industry leaders choose to Build with Devin
        </h2>
        {/* Optional: Add a short sentence or two here if there's relevant subtext */}
        {/* <p className="text-lg text-gray-300 mb-10 max-w-xl mx-auto">
          Join the growing number of innovative companies transforming their development processes with Devin.
        </p> */}
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-8 rounded-lg text-lg transition duration-300 ease-in-out transform hover:scale-105"
          onClick={() => {
            // Placeholder action, replace with actual link or modal opening
            alert('Navigating to customer stories... (placeholder)');
          }}
        >
          Hear from our customers
        </button>
      </div>
    </section>
  );
};

export default IndustryLeadersSection;
