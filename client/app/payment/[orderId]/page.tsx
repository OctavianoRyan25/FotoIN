"use client";

import { Suspense, useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import Image from "next/image";
import { CheckCircle, Loader2, ArrowLeft } from "lucide-react";
import { toast } from "sonner";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

function PaymentPageContent() {
  const router = useRouter();
  const { orderId } = useParams<{ orderId: string }>();
  const searchParams = useSearchParams();
  const qrCodeUrl = searchParams.get("qr");

  const [loading, setLoading] = useState(false);
  const [paymentSuccess, setPaymentSuccess] = useState(false);
  const [orderData, setOrderData] = useState<any>(null);
  const buyerEmail = searchParams.get("email") ?? "";

  useEffect(() => {
    if (!qrCodeUrl) {
      toast.error("QR Code tidak ditemukan");
      router.push("/");
    }
  }, [qrCodeUrl, router]);

  const checkPaymentStatus = async () => {
    try {
      setLoading(true);
      const query = buyerEmail
        ? `?email=${encodeURIComponent(buyerEmail)}`
        : "";
      const res = await fetch(`${API_URL}/public/orders/${orderId}${query}`);
      const data = await res.json();

      if (data.data?.payment_status === "paid") {
        setPaymentSuccess(true);
        setOrderData(data.data);
        toast.success(
          "Pembayaran berhasil! Cek email Anda untuk link download.",
        );
      } else {
        toast.info("Pembayaran belum terdeteksi. Silakan coba lagi.");
      }
    } catch (error) {
      toast.error("Gagal memeriksa status pembayaran");
    } finally {
      setLoading(false);
    }
  };

  if (paymentSuccess) {
    return (
      <div className="min-h-screen bg-linear-to-br from-violet-50 to-purple-50 flex items-center justify-center p-4">
        <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8 text-center">
          <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <CheckCircle className="w-10 h-10 text-green-600" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Pembayaran Berhasil!
          </h1>
          <p className="text-gray-600 mb-6">
            Link download foto telah dikirim ke email Anda. Cek folder inbox
            atau spam.
          </p>
          <div className="bg-gray-50 rounded-lg p-4 mb-6 text-left">
            <div className="flex justify-between text-sm mb-2">
              <span className="text-gray-600">Order ID</span>
              <span className="font-mono text-xs">{orderId}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Total Foto</span>
              <span className="font-semibold">
                {orderData?.order_items?.length || 0} foto
              </span>
            </div>
          </div>
          <button
            onClick={() => router.push("/")}
            className="w-full bg-violet-600 text-white py-3 rounded-lg font-medium hover:bg-violet-700 transition"
          >
            Kembali ke Beranda
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-linear-to-br from-violet-50 to-purple-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8">
        <button
          onClick={() => router.back()}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6 transition"
        >
          <ArrowLeft className="w-4 h-4" />
          <span className="text-sm">Kembali</span>
        </button>

        <div className="text-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Scan QR Code untuk Bayar
          </h1>
          <p className="text-gray-600 text-sm">
            Gunakan aplikasi e-wallet Anda untuk scan QR code di bawah
          </p>
        </div>

        <div className="bg-white border-4 border-violet-100 rounded-2xl p-6 mb-6">
          {qrCodeUrl && (
            <div className="relative w-full aspect-square">
              <Image
                src={qrCodeUrl}
                alt="QR Code Pembayaran"
                fill
                className="object-contain"
                unoptimized
              />
            </div>
          )}
        </div>

        <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 mb-6">
          <p className="text-sm text-amber-800">
            <strong>QR Code berlaku 30 menit.</strong> Setelah melakukan
            pembayaran, klik tombol di bawah untuk mengecek status.
          </p>
        </div>

        <button
          onClick={checkPaymentStatus}
          disabled={loading}
          className="w-full bg-violet-600 text-white py-3 rounded-lg font-medium hover:bg-violet-700 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {loading ? (
            <>
              <Loader2 className="w-5 h-5 animate-spin" />
              Mengecek Status...
            </>
          ) : (
            "Saya Sudah Bayar"
          )}
        </button>

        <p className="text-xs text-gray-500 text-center mt-4">
          Order ID: <span className="font-mono">{orderId}</span>
        </p>
      </div>
    </div>
  );
}

export default function PaymentPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          Memuat pembayaran...
        </div>
      }
    >
      <PaymentPageContent />
    </Suspense>
  );
}
