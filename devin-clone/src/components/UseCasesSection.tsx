import React from 'react';

interface UseCaseCardProps {
  imagePlaceholder: string;
  title: string;
  description: string;
}

const UseCaseCard: React.FC<UseCaseCardProps> = ({ imagePlaceholder, title, description }) => {
  return (
    <div className="bg-white p-6 rounded-lg shadow-lg flex flex-col items-center text-center">
      <div className="w-full h-40 bg-gray-300 rounded-md flex items-center justify-center mb-4">
        <p className="text-gray-500">{imagePlaceholder}</p>
      </div>
      <h3 className="text-xl font-semibold text-gray-800 mb-2">{title}</h3>
      <p className="text-gray-600 text-sm">{description}</p>
    </div>
  );
};

const UseCasesSection = () => {
  const useCases = [
    {
      imagePlaceholder: "Migration/Refactor Icon",
      title: "Code Migration + Refactors",
      description: "Automate large-scale code migrations and refactor legacy systems with AI precision."
    },
    {
      imagePlaceholder: "Data Science Icon",
      title: "Data Science",
      description: "Accelerate data analysis, model training, and insight generation for your data science workflows."
    },
    {
      imagePlaceholder: "Triage Icon",
      title: "Issue Triage",
      description: "Automatically categorize, prioritize, and assign incoming issues and bug reports."
    },
    {
      imagePlaceholder: "CI/CD Icon",
      title: "CI/CD",
      description: "Optimize and manage continuous integration and deployment pipelines efficiently."
    },
    {
      imagePlaceholder: "Repo Setup Icon",
      title: "Repository Setup",
      description: "Standardize and automate the setup of new repositories with best practices."
    },
    {
      imagePlaceholder: "General Task Icon",
      title: "And much more...",
      description: "Devins are versatile and can be trained for a wide variety of software engineering tasks."
    }
  ];

  return (
    <section className="py-12 md:py-24 bg-white">
      <div className="container mx-auto px-6">
        <h2 className="text-3xl md:text-4xl font-bold text-center text-gray-800 mb-16">
          Devins can work tirelessly and in parallel
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {useCases.map((useCase, index) => (
            <UseCaseCard
              key={index}
              imagePlaceholder={useCase.imagePlaceholder}
              title={useCase.title}
              description={useCase.description}
            />
          ))}
        </div>
      </div>
    </section>
  );
};

export default UseCasesSection;
