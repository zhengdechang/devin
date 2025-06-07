import React from 'react';

const Hero = () => {
  return (
    <section className="relative bg-gray-800 text-white py-20 md:py-32">
      {/* Background Image Placeholder */}
      {/* Replace with <Image /> component from Next.js for optimized images if needed */}
      <div
        className="absolute inset-0 bg-cover bg-center opacity-50"
        style={{ backgroundImage: "url('https://via.placeholder.com/1920x1080/000022/FFFFFF?Text=Hero+Background')" }}
      ></div>

      {/* Overlay to darken the background image for better text readability */}
      <div className="absolute inset-0 bg-black opacity-50"></div>

      <div className="container mx-auto px-6 text-center relative z-10">
        <h1 className="text-4xl md:text-6xl font-bold mb-6">
          Revolutionize Your Workflow
        </h1>
        <p className="text-lg md:text-xl mb-8 max-w-2xl mx-auto">
          Empower your team with cutting-edge AI solutions. Drive efficiency, foster innovation, and achieve unparalleled results.
        </p>
        <button className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-8 rounded-lg text-lg">
          Get Started Today
        </button>
      </div>
    </section>
  );
};

export default Hero;
