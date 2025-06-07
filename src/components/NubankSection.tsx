import React from 'react';

const NubankSection = () => {
  return (
    <section className="py-12 md:py-24 bg-white">
      <div className="container mx-auto px-6">
        {/* Section Title */}
        <h2 className="text-3xl md:text-4xl font-bold text-center text-gray-800 mb-6">
          How Nubank refactors with Large Context Models
        </h2>
        {/* Section Description */}
        <p className="text-lg text-gray-600 text-center max-w-3xl mx-auto mb-12">
          Discover how Nubank leverages cutting-edge AI to streamline their refactoring process, achieving significant gains in efficiency and cost savings. Our case study highlights the transformative power of Large Context Models in a real-world enterprise scenario.
        </p>

        <div className="flex flex-col md:flex-row md:space-x-8 items-center">
          {/* Video Placeholder */}
          <div className="w-full md:w-1/2 mb-8 md:mb-0">
            <div className="aspect-video bg-gray-200 rounded-lg flex items-center justify-center">
              {/* In a real scenario, you'd embed the Vimeo player here */}
              <p className="text-gray-500">Vimeo Video Placeholder (16:9)</p>
            </div>
            <p className="text-sm text-gray-500 mt-2 text-center">
              Watch the full story on how Nubank transformed their development lifecycle.
            </p>
          </div>

          {/* Statistics and Details */}
          <div className="w-full md:w-1/2">
            <h3 className="text-2xl font-semibold text-gray-700 mb-4">Key Outcomes:</h3>
            <ul className="space-y-4">
              <li className="bg-gray-50 p-4 rounded-lg shadow">
                <p className="text-xl font-bold text-blue-600">8x Engineering Time Efficiency Gain</p>
                <p className="text-gray-600">Reduced time spent on manual refactoring tasks by a factor of eight.</p>
              </li>
              <li className="bg-gray-50 p-4 rounded-lg shadow">
                <p className="text-xl font-bold text-blue-600">20x Cost Savings</p>
                <p className="text-gray-600">Significant reduction in operational costs related to code maintenance and development.</p>
              </li>
              <li className="bg-gray-50 p-4 rounded-lg shadow">
                <p className="text-xl font-bold text-blue-600">Improved Code Quality</p>
                <p className="text-gray-600">Enhanced code consistency and maintainability through AI-assisted refactoring.</p>
              </li>
            </ul>
            <a
              href="/case-studies/nubank" // Placeholder link
              className="mt-8 inline-block bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-6 rounded-lg text-center"
            >
              Read Full Case Study
            </a>
          </div>
        </div>
      </div>
    </section>
  );
};

export default NubankSection;
