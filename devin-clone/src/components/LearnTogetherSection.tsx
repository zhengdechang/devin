import React from 'react';

interface LearnFeatureCardProps {
  imagePlaceholder: string;
  title: string;
  description: string;
}

const LearnFeatureCard: React.FC<LearnFeatureCardProps> = ({ imagePlaceholder, title, description }) => {
  return (
    <div className="flex-1 bg-white p-6 rounded-lg shadow-lg text-center md:text-left">
      <div className="w-full h-48 bg-gray-300 rounded-md flex items-center justify-center mb-4 md:float-left md:w-1/3 md:h-auto md:mr-6">
        <p className="text-gray-500">{imagePlaceholder}</p>
      </div>
      <div>
        <h3 className="text-xl font-semibold text-gray-800 mb-2">{title}</h3>
        <p className="text-gray-600">{description}</p>
      </div>
    </div>
  );
};


const LearnTogetherSection = () => {
  const learnFeatures = [
    {
      imagePlaceholder: "Codebase Learning Icon",
      title: "Devin learns your codebase & picks up tribal knowledge",
      description: "Continuously improves by observing your development practices and learning from your existing code and documentation."
    },
    {
      imagePlaceholder: "Mobile Coding Icon",
      title: "Code on the go",
      description: "Access Devin's capabilities from anywhere, allowing you to manage tasks and review code remotely."
    },
    {
      imagePlaceholder: "Integrated Tools Icon",
      title: "Use Devin's editor, shell and browser",
      description: "Leverage Devin's built-in tools for a seamless development experience, from coding to testing and debugging."
    }
  ];

  return (
    <section className="py-12 md:py-24 bg-gray-100">
      <div className="container mx-auto px-6">
        <h2 className="text-3xl md:text-4xl font-bold text-center text-gray-800 mb-16">
          Learn & work together
        </h2>
        <div className="space-y-12">
          {learnFeatures.map((feature, index) => (
            <div key={index} className="md:flex items-center bg-white p-6 rounded-lg shadow-lg">
              <div className="md:w-1/3 h-48 md:h-auto bg-gray-300 rounded-md flex items-center justify-center mb-4 md:mb-0 md:mr-6">
                 <p className="text-gray-500 text-sm">{feature.imagePlaceholder}</p>
              </div>
              <div className="md:w-2/3">
                <h3 className="text-2xl font-semibold text-gray-800 mb-3">{feature.title}</h3>
                <p className="text-gray-600">{feature.description}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default LearnTogetherSection;
