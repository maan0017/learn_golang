"use client";

import { ChangeEvent, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Image from "next/image";

type OptimizedImageProps = {
  publicId: string;
};

function OptimizedImage({ publicId }: OptimizedImageProps) {
  return (
    <Image
      src={publicId}
      alt="Uploaded"
      className="rounded-xl shadow-md mt-4 w-72"
    />
  );
}

export default function CloudinaryUploadDemo() {
  const [file, setFile] = useState<File | null>(null);
  const [imageUrl, setImageUrl] = useState<string>("");

  const handleUpload = async () => {
    if (!file) return;
    console.log("sending file...");

    const formData = new FormData();
    formData.append("file", file);
    console.log("Form file created...");

    // 👉 call your Go backend /upload endpoint
    const res = await fetch("http://localhost:8080/upload", {
      method: "POST",
      body: formData,
    });

    const data = await res.json();
    setImageUrl(data.url);
    console.log(data.url);
    console.log("file uploaded...");
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-full max-w-md p-6 shadow-lg rounded-2xl">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">
            Cloudinary Upload Demo
          </CardTitle>
        </CardHeader>

        <CardContent className="space-y-4">
          <Input
            type="file"
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setFile(e.target.files?.[0] || null)
            }
            className="cursor-pointer"
          />

          <Button
            type="button"
            onClick={handleUpload}
            disabled={!file}
            className="w-full cursor-pointer disabled:cursor-not-allowed"
          >
            Upload
          </Button>

          {imageUrl && (
            <div className="text-center">
              <h3 className="text-lg font-medium">Uploaded Image:</h3>
              <OptimizedImage publicId={imageUrl} />
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
