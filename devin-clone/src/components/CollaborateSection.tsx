import React from 'react';

interface FeatureCardProps {
  imagePlaceholder: string; // Simple text for placeholder color/content
  title: string;
  description: string;
}

const FeatureCard: React.FC<FeatureCardProps> = ({ imagePlaceholder, title, description }) => {
  return (
    <div className="flex-1 bg-gray-50 p-6 rounded-lg shadow-lg text-center">
      <div className="w-full h-48 bg-gray-300 rounded-md flex items-center justify-center mb-4">
        <p className="text-gray-500">{imagePlaceholder}</p>
      </div>
      <h3 className="text-xl font-semibold text-gray-800 mb-2">{title}</h3>
      <p className="text-gray-600">{description}</p>
    </div>
  );
};

const CollaborateSection = () => {
  const features = [
    {
      imagePlaceholder: "Shell Icon/Graphic",
      title: "Its Own Shell",
      description: "Devin is equipped with its own sandboxed shell environment, allowing it to execute commands and scripts for your projects seamlessly."
    },
    {
      imagePlaceholder: "Editor Icon/Graphic",
      title: "Its Own Code Editor",
      description: "Featuring a fully integrated code editor, Devin can write, review, and debug code autonomously or with your guidance."
    },
    {
      imagePlaceholder: "Browser Icon/Graphic",
      title: "Its Own Browser",
      description: "Devin utilizes its own browser to interact with web interfaces, gather information, and perform tasks just like a human would."
    }
  ];

  return (
    <section className="py-12 md:py-24 bg-gray-100">
      <div className="container mx-auto px-6">
        <h2 className="text-3xl md:text-4xl font-bold text-center text-gray-800 mb-16">
          Built to collaborate with you
        </h2>
        <div className="flex flex-col md:flex-row gap-8 justify-center items-stretch">
          {features.map((feature, index) => (
            <FeatureCard
              key={index}
              imagePlaceholder={feature.imagePlaceholder}
              title={feature.title}
              description={feature.description}
            />
          ))}
        </div>
      </div>
    </section>
  );
};

export default CollaborateSection;
